// Command seed fills a dev database with a sample user and a handful of
// sample posts (with generated photos uploaded to S3), so the blog has
// something to look at locally.
//
// With the docker compose stack running:
//
//	docker compose exec web go run ./cmd/seed
//
// It refuses to run against a database that already has posts unless
// -force (add alongside) or -reset (delete every post first) is passed.
package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"flag"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	mrand "math/rand/v2"
	"os"
	"time"

	"doppler/internal/db"
	"doppler/internal/services"
	"doppler/internal/storage"

	"github.com/joho/godotenv"
)

type samplePost struct {
	title   string
	content string
	daysAgo float64
	photos  []palette
}

var samplePosts = []samplePost{
	{
		title:   "things my houseplant has survived",
		daysAgo: 230,
		content: `<p>A partial list, in no particular order:</p>
<ul><li>two moves</li><li>a week of vacation with zero water</li><li>being knocked off the shelf by a cat that is not mine</li><li>me reading an article about overwatering and then overwatering it</li></ul>
<blockquote>it is thriving. i don't know why. i've stopped asking.</blockquote>`,
	},
	{
		title:   "weekend contact sheet",
		daysAgo: 70,
		content: `<p>Took the camera out with no plan at all. Five keepers out of about two hundred, which feels like a good ratio for a Saturday.</p>`,
		photos:  []palette{dusk, meadow, harbor, lavender, ember},
	},
	{
		title:   "notes from the server closet",
		daysAgo: 41,
		content: `<p>Photos on this blog now live in <strong>Garage</strong> and everything else lives in SQLite. The whole upload path is roughly:</p>
<pre>browser → echo → garage → back to your browser</pre>
<p>Is it overkill for a blog that gets posted to once a month? Absolutely. Did I have fun? Also absolutely. The rest of the projects are over on <a href="https://pinwheel.fan">pinwheel.fan</a>.</p>`,
		photos: []palette{night, harbor},
	},
	{
		title:   "a very scientific ranking of clouds",
		daysAgo: 23,
		content: `<p>After years of research (looking up), the official results are in.</p>
<ol><li><strong>Cumulus.</strong> Classic. Looks like a sheep. No notes.</li><li><strong>Mammatus.</strong> Unsettling, but in a fun way.</li><li><strong>Whatever that one was on Tuesday.</strong> Looked exactly like a dog wearing a hat.</li></ol>
<p>Cirrus did not place. Cirrus knows what it did.</p>`,
		photos: []palette{meadow, lavender, dusk},
	},
	{
		title:   "the sun did a thing tonight",
		daysAgo: 9,
		content: `<p>Walked outside to take the recycling out and the whole sky was on fire. The recycling did not get taken out.</p>`,
		photos:  []palette{ember},
	},
	{
		title:   "hello, world (again)",
		daysAgo: 2,
		content: `<p>New coat of paint around here. Same blog, fewer moving parts.</p>
<p>Expect photos, half-finished thoughts, and the occasional ranking of things that do not need to be ranked.</p>`,
	},
}

func main() {
	username := flag.String("user", "pinwheel", "username for the sample author (created if missing)")
	password := flag.String("password", "pinwheel", "password for a newly created sample author")
	force := flag.Bool("force", false, "add sample posts even if the database already has posts")
	reset := flag.Bool("reset", false, "delete ALL existing posts and pictures before seeding (dev only)")
	flag.Parse()

	_ = godotenv.Load()
	// In the compose dev container, S3 credentials are written to a shared file by garage-init.
	if os.Getenv("S3_ACCESS_KEY_ID") == "" {
		_ = godotenv.Load("/shared/s3-credentials.env")
	}

	conn := db.Connect()
	if err := db.RunMigrations(conn); err != nil {
		log.Fatalf("migrations: %v", err)
	}

	if *reset {
		if _, err := conn.Exec("DELETE FROM picture; DELETE FROM post;"); err != nil {
			log.Fatalf("reset: %v", err)
		}
		log.Printf("deleted all existing posts")
	}

	var existing int
	if err := conn.QueryRow("SELECT count(*) FROM post").Scan(&existing); err != nil {
		log.Fatalf("counting posts: %v", err)
	}
	if existing > 0 && !*force {
		log.Fatalf("database already has %d posts; pass -force to add samples anyway", existing)
	}

	user, err := services.GetUserByUsername(conn, *username)
	if err != nil || user == nil {
		user = services.CreateUser(conn, *username, *password, *username+"@example.com")
		log.Printf("created user %q (password %q)", *username, *password)
	}

	bucket := os.Getenv("S3_BUCKET")
	client := storage.NewGarageClient()

	for _, sample := range samplePosts {
		post := services.CreatePost(conn, user.ID, sample.title, sample.content)
		// Whole days back plus a random time of day, so posts don't all share a clock time.
		created := time.Now().UTC().Add(-time.Duration(sample.daysAgo*float64(24*time.Hour)) - time.Duration(mrand.IntN(12*60))*time.Minute)
		if _, err := conn.Exec("UPDATE post SET created = ?, modified = ? WHERE id = ?", created, created, post.ID); err != nil {
			log.Fatalf("dating post %d: %v", post.ID, err)
		}
		for _, p := range sample.photos {
			if err := addPhoto(conn, client, bucket, post.ID, p); err != nil {
				log.Fatalf("photo for %q: %v", sample.title, err)
			}
		}
		log.Printf("seeded %q (%d photos)", sample.title, len(sample.photos))
	}
}

func addPhoto(conn *sql.DB, client *storage.GarageClient, bucket string, postID int, p palette) error {
	name := make([]byte, 16)
	if _, err := rand.Read(name); err != nil {
		return err
	}
	filename := hex.EncodeToString(name)

	f, err := os.CreateTemp("", "doppler-seed-*.png")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	if err := png.Encode(f, paint(1200, 900, p)); err != nil {
		f.Close()
		return err
	}
	f.Close()

	client.StoreObject(bucket, filename, f.Name())
	_, err = conn.Exec("INSERT INTO picture (post_id, filename) VALUES (?, ?)", postID, filename)
	return err
}

// palette describes a flat little landscape: sky, sun, far hills, near hills, ground.
type palette struct {
	sky, sun, far, near, ground color.RGBA
	stars                       bool
}

func hexColor(s string) color.RGBA {
	b, _ := hex.DecodeString(s)
	return color.RGBA{b[0], b[1], b[2], 255}
}

var (
	dusk     = palette{hexColor("f4a261"), hexColor("e76f51"), hexColor("2a9d8f"), hexColor("264653"), hexColor("1d3557"), false}
	meadow   = palette{hexColor("bde0fe"), hexColor("fff3b0"), hexColor("a3b18a"), hexColor("588157"), hexColor("344e41"), false}
	harbor   = palette{hexColor("90e0ef"), hexColor("ffffff"), hexColor("0096c7"), hexColor("023e8a"), hexColor("03045e"), false}
	lavender = palette{hexColor("cdb4db"), hexColor("ffafcc"), hexColor("a2d2ff"), hexColor("577590"), hexColor("22223b"), false}
	ember    = palette{hexColor("ffd166"), hexColor("ef476f"), hexColor("118ab2"), hexColor("073b4c"), hexColor("06141b"), false}
	night    = palette{hexColor("141b2d"), hexColor("f1f1e6"), hexColor("2b3a55"), hexColor("1b2537"), hexColor("0d121c"), true}
)

// paint draws a flat landscape with randomly shaped ridges, so every photo differs a little.
func paint(w, h int, p palette) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	r := mrand.New(mrand.NewPCG(mrand.Uint64(), mrand.Uint64()))

	ridge := func(base, amp float64) []float64 {
		ys := make([]float64, w)
		phase1, phase2 := r.Float64()*6, r.Float64()*6
		f1, f2 := 2+r.Float64()*3, 5+r.Float64()*5
		for x := range ys {
			t := float64(x) / float64(w)
			// Triangle waves make pointy hills; the second adds wobble.
			tri := 2*math.Abs(math.Mod(t*f1+phase1, 1)-0.5) - 0.5
			ys[x] = base + amp*tri + amp*0.25*math.Sin(t*f2*math.Pi*2+phase2)
		}
		return ys
	}
	far := ridge(float64(h)*0.55, float64(h)*0.18)
	near := ridge(float64(h)*0.70, float64(h)*0.12)
	sunX, sunY := float64(w)*(0.2+r.Float64()*0.6), float64(h)*(0.18+r.Float64()*0.15)
	sunR := float64(h) * (0.07 + r.Float64()*0.05)

	for y := range h {
		for x := range w {
			fx, fy := float64(x), float64(y)
			c := p.sky
			if p.stars && r.IntN(900) == 0 {
				c = p.sun
			}
			if (fx-sunX)*(fx-sunX)+(fy-sunY)*(fy-sunY) < sunR*sunR {
				c = p.sun
			}
			if fy > far[x] {
				c = p.far
			}
			if fy > near[x] {
				c = p.near
			}
			if fy > float64(h)*0.86+float64(h)*0.03*math.Sin(fx/float64(w)*math.Pi*2) {
				c = p.ground
			}
			img.SetRGBA(x, y, c)
		}
	}
	return img
}
