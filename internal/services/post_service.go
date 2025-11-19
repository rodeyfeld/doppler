package services

import (
	"crypto/rand"
	"database/sql"
	"doppler/internal/models"
	"doppler/internal/storage"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
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

// GetPostByID retrieves a single post with its pictures
func GetPostByID(db *sql.DB, postID int) (models.Post, error) {
	// Query the post
	postQuery := `
		SELECT id, user_id, title, content, created, modified
		FROM post
		WHERE id = ?
	`

	var post models.Post
	err := db.QueryRow(postQuery, postID).Scan(
		&post.ID, &post.UserID, &post.Title, &post.Content, &post.Created, &post.Modified,
	)
	if err != nil {
		return models.Post{}, fmt.Errorf("failed to get post: %v", err)
	}

	post.PictureURLs = []string{} // Initialize empty slice

	// Query pictures for this post
	pictureQuery := `
		SELECT filename
		FROM picture
		WHERE post_id = ?
		ORDER BY id ASC
	`

	rows, err := db.Query(pictureQuery, postID)
	if err != nil {
		log.Printf("Picture query failed: %v", err)
		return post, nil // Return post without pictures
	}
	defer rows.Close()

	for rows.Next() {
		var filename string
		if err := rows.Scan(&filename); err != nil {
			log.Printf("Failed scanning picture: %v", err)
			continue
		}
		pictureURL := fmt.Sprintf("/doppler/images/%s", filename)
		post.PictureURLs = append(post.PictureURLs, pictureURL)
	}

	return post, nil
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

// validateImage validates that the uploaded file is a valid image
func validateImage(formFile *multipart.FileHeader) error {
	// 1. Check file size (max 10MB)
	const maxFileSize = 10 * 1024 * 1024 // 10MB
	if formFile.Size > maxFileSize {
		return errors.New("file too large: maximum size is 10MB")
	}

	// 2. Check MIME type from header
	contentType := formFile.Header.Get("Content-Type")
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}
	if !allowedTypes[contentType] {
		return fmt.Errorf("invalid content type: %s (allowed: jpeg, png, gif, webp)", contentType)
	}

	// 3. Verify actual image format by decoding
	src, err := formFile.Open()
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer src.Close()

	// Try to decode as image - this validates it's actually an image
	_, format, err := image.DecodeConfig(src)
	if err != nil {
		return fmt.Errorf("invalid image file: %v", err)
	}

	// Validate format matches one we support
	allowedFormats := map[string]bool{
		"jpeg": true,
		"png":  true,
		"gif":  true,
	}
	if !allowedFormats[format] {
		return fmt.Errorf("unsupported image format: %s", format)
	}

	log.Printf("Image validated: format=%s, size=%d bytes", format, formFile.Size)
	return nil
}

// createTempImage creates a temp file and copies uploaded file data to it
func createTempImage(formFile *multipart.FileHeader) (string, error) {
	// Open source file from form upload
	src, err := formFile.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer src.Close()

	// Create temporary file in system tmp directory
	dst, err := os.CreateTemp(os.TempDir(), "doppler-upload-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %v", err)
	}
	defer dst.Close()

	// Copy uploaded file to temp file
	if _, err = io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to copy file: %v", err)
	}

	return dst.Name(), nil
}

func CreatePicture(db *sql.DB, postID int, formFile *multipart.FileHeader) (models.Picture, error) {
	// Validate image before processing
	if err := validateImage(formFile); err != nil {
		return models.Picture{}, fmt.Errorf("image validation failed: %v", err)
	}

	s3PhotoBucket := os.Getenv("S3_BUCKET")
	hexLength := 32
	bytes := make([]byte, hexLength/2) // n hex chars requires n/2 bytes
	if _, err := rand.Read(bytes); err != nil {
		return models.Picture{}, fmt.Errorf("failed to generate random bytes: %v", err)
	}

	filename := hex.EncodeToString(bytes)
	filepath, err := createTempImage(formFile)
	if err != nil {
		return models.Picture{}, err
	}
	defer os.Remove(filepath) // Clean up temp file after function exits

	uploadObject(s3PhotoBucket, filename, filepath)

	var query, dbErr = db.Prepare("INSERT INTO picture (post_id, filename) VALUES (?, ?) RETURNING id, post_id, filename")
	if dbErr != nil {
		return models.Picture{}, fmt.Errorf("failed to prepare query: %v", dbErr)
	}

	picture := models.Picture{}
	dbErr = query.QueryRow(postID, filename).Scan(&picture.ID, &picture.PostID, &picture.Filename)
	if dbErr != nil {
		return models.Picture{}, fmt.Errorf("query failed: %v", dbErr)
	}

	return picture, nil
}

func uploadObject(bucketName string, filename string, filepath string) {
	client := storage.NewGarageClient()
	client.StoreObject(bucketName, filename, filepath)
}

// ValidateFilename validates that the filename is safe (no path traversal)
func ValidateFilename(filename string) error {
	// Check for empty filename
	if filename == "" {
		return errors.New("filename cannot be empty")
	}

	// Check for path traversal attempts
	if len(filename) > 100 {
		return errors.New("filename too long")
	}

	// Only allow hexadecimal characters (our generated filenames are hex)
	for _, char := range filename {
		if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f') || (char >= 'A' && char <= 'F')) {
			return fmt.Errorf("invalid character in filename: %c", char)
		}
	}

	// Our filenames are exactly 32 hex characters
	if len(filename) != 32 {
		return fmt.Errorf("invalid filename length: expected 32, got %d", len(filename))
	}

	return nil
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
