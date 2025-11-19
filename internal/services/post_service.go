package services

import (
	"crypto/rand"
	"database/sql"
	"doppler/internal/models"
	"doppler/internal/storage"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func GetPosts(db *sql.DB) []models.Post {
	// Query all posts
	postQuery := `
		SELECT id, user_id, title, content, created, modified
		FROM post
		ORDER BY created DESC
	`

	rows, err := db.Query(postQuery)
	if err != nil {
		log.Printf("Query failed: %v", err)
		return []models.Post{}
	}
	defer rows.Close()

	posts := []models.Post{}
	postIDs := []int{}

	for rows.Next() {
		var p models.Post
		err = rows.Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &p.Created, &p.Modified)
		if err != nil {
			log.Printf("Failed scanning to Post: %v", err)
			continue
		}
		p.PictureURLs = []string{} // Initialize empty slice
		posts = append(posts, p)
		postIDs = append(postIDs, p.ID)
	}

	// If no posts, return early
	if len(posts) == 0 {
		return posts
	}

	// Query all pictures for these posts
	pictureQuery := `
		SELECT post_id, filename
		FROM picture
		WHERE post_id IN (` + placeholders(len(postIDs)) + `)
		ORDER BY id ASC
	`

	// Convert postIDs to interface slice for query
	args := make([]interface{}, len(postIDs))
	for i, id := range postIDs {
		args[i] = id
	}

	pictureRows, err := db.Query(pictureQuery, args...)
	if err != nil {
		log.Printf("Picture query failed: %v", err)
		return posts
	}
	defer pictureRows.Close()

	// Map pictures to posts
	postMap := make(map[int]*models.Post)
	for i := range posts {
		postMap[posts[i].ID] = &posts[i]
	}

	for pictureRows.Next() {
		var postID int
		var filename string

		err = pictureRows.Scan(&postID, &filename)
		if err != nil {
			log.Printf("Failed scanning picture: %v", err)
			continue
		}

		// Generate backend proxy URL for picture
		pictureURL := fmt.Sprintf("/doppler/images/%s", filename)

		// Add URL to corresponding post
		if post, ok := postMap[postID]; ok {
			post.PictureURLs = append(post.PictureURLs, pictureURL)
		}
	}

	return posts
}

// placeholders generates SQL placeholder string like "?,?,?"
func placeholders(n int) string {
	if n <= 0 {
		return ""
	}
	result := "?"
	for i := 1; i < n; i++ {
		result += ",?"
	}
	return result
}

func CreatePost(db *sql.DB, userID int, title string, content string) models.Post {

	var query, err = db.Prepare("INSERT INTO post (user_id, title, content) VALUES (?, ?, ?) RETURNING id, user_id, title, content")
	if err != nil {
		log.Panic(err)
	}

	post := models.Post{}
	err = query.QueryRow(userID, title, content).Scan(&post.ID, &post.UserID, &post.Title, &post.Content)
	if err != nil {
		log.Panicf("Query failed: %v", err)
	}

	return post
}

// createTempImage creates a temp file and copies uploaded file data to it
func createTempImage(formFile *multipart.FileHeader) string {
	// Open source file from form upload
	src, err := formFile.Open()
	if err != nil {
		log.Panicf("Failed to open file: %v", err)
	}
	defer src.Close()

	// Create temporary file in system tmp directory
	dst, err := os.CreateTemp(os.TempDir(), "doppler-upload-*")
	if err != nil {
		log.Panicf("Failed to create temp file: %v", err)
	}
	defer dst.Close()

	// Copy uploaded file to temp file
	if _, err = io.Copy(dst, src); err != nil {
		log.Panicf("Failed to copy file: %v", err)
	}

	return dst.Name()
}

func CreatePicture(db *sql.DB, postID int, formFile *multipart.FileHeader) models.Picture {
	s3PhotoBucket := os.Getenv("S3_BUCKET")
	hexLength := 32
	bytes := make([]byte, hexLength/2) // n hex chars requires n/2 bytes
	if _, err := rand.Read(bytes); err != nil {
		log.Panicf("Failed to generate random bytes")
	}

	filename := hex.EncodeToString(bytes)
	filepath := createTempImage(formFile)
	defer os.Remove(filepath) // Clean up temp file after function exits
	uploadObject(s3PhotoBucket, filename, filepath)

	var query, err = db.Prepare("INSERT INTO picture (post_id, filename) VALUES (?, ?) RETURNING id, post_id, filename")
	if err != nil {
		log.Panic(err)
	}

	picture := models.Picture{}
	err = query.QueryRow(postID, filename).Scan(&picture.ID, &picture.PostID, &picture.Filename)
	if err != nil {
		log.Panicf("Query failed: %v", err)
	}
	return picture
}

func uploadObject(bucketName string, filename string, filepath string) {
	client := storage.NewGarageClient()
	client.StoreObject(bucketName, filename, filepath)
}

// GetPictureStream fetches an image from S3 and returns a reader for streaming
func GetPictureStream(filename string) (io.ReadCloser, string, error) {
	s3PhotoBucket := os.Getenv("S3_BUCKET")
	client := storage.NewGarageClient()

	reader, contentType, err := client.GetObjectStream(s3PhotoBucket, filename)
	if err != nil {
		return nil, "", err
	}

	return reader, contentType, nil
}
