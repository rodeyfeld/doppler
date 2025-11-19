package handlers

import (
	"doppler/internal/components/post"
	"doppler/internal/components/shared"
	"doppler/internal/models"
	"doppler/internal/server"
	"doppler/internal/services"
	"fmt"
	"log"
	"strconv"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/microcosm-cc/bluemonday"
)

type PostHandler struct {
	server *server.DopplerServer
}

func NewPostHandler(s *server.DopplerServer) *PostHandler {
	return &PostHandler{server: s}
}

func (h *PostHandler) Index(c echo.Context) error {
	posts := services.GetPosts(h.server.DB)
	sess, err := session.Get("auth-session", c)
	if err != nil {
		log.Printf("Failed to get auth session")
		return err
	}

	var user *models.User
	// Check if user is logged in
	if userID, ok := sess.Values["userID"].(int); ok {
		user, err = services.GetUserByID(h.server.DB, userID)
		if err != nil {
			log.Printf("Failed to get user by ID: %v", err)
		}
	}

	cmp := post.PostIndex(post.ListPosts(posts), user)
	return renderView(c, cmp)
}

func (h *PostHandler) Create(c echo.Context) error {
	sess, err := session.Get("auth-session", c)
	if err != nil {
		log.Printf("Failed to get auth session: %v", err)
		return err
	}

	// Check if user is logged in
	userID, ok := sess.Values["userID"].(int)
	if !ok {
		log.Printf("User not authenticated, cannot create post")
		cmp := shared.AuthRequired()
		return renderView(c, cmp)
	}

	// Get form values
	rawTitle := c.FormValue("title")
	rawContent := c.FormValue("content")
	imageFormFile, err := c.FormFile("image-content")
	if err != nil {
		log.Printf("No image file provided: %v", err)
		return c.String(400, "Image file is required")
	}

	// Sanitize inputs using bluemonday
	// UGCPolicy allows safe HTML tags (for Quill content) while stripping dangerous elements
	p := bluemonday.UGCPolicy()

	// Sanitize title (strip all HTML tags)
	title := bluemonday.StrictPolicy().Sanitize(rawTitle)

	// Sanitize content (allow safe HTML tags for rich text)
	content := p.Sanitize(rawContent)
	log.Printf("Creating post - Title: %s, Content length: %d", title, len(content))
	log.Printf("Content preview (sanitized): %s", content)

	createdPost := services.CreatePost(h.server.DB, userID, title, content)
	createdImage, err := services.CreatePicture(h.server.DB, createdPost.ID, imageFormFile)
	if err != nil {
		log.Printf("Failed to create picture: %v", err)
		return c.String(400, fmt.Sprintf("Image upload failed: %v", err))
	}

	cmp := post.PostSuccess(createdPost, createdImage)
	return renderView(c, cmp)
}

func (h *PostHandler) UserInfo(c echo.Context) error {
	p := c.Param("id")
	log.Printf("User info requested for ID: %s", p)
	id, err := strconv.Atoi(p)
	if err != nil {

		return err
	}
	user, err := services.GetUserByID(h.server.DB, id)
	if err != nil {
		return err
	}
	cmp := post.PostUserInfo(user)
	return renderView(c, cmp)
}

// GetImage serves images from S3 storage through the backend as a proxy
// This avoids exposing S3/Garage endpoints publicly
func (h *PostHandler) GetImage(c echo.Context) error {
	filename := c.Param("filename")
	if filename == "" {
		return c.String(400, "Filename is required")
	}

	// Validate filename to prevent path traversal attacks
	if err := services.ValidateFilename(filename); err != nil {
		log.Printf("Invalid filename requested: %s, error: %v", filename, err)
		return c.String(400, "Invalid filename")
	}

	// Stream the image from S3
	reader, contentType, err := services.GetPictureStream(filename)
	if err != nil {
		log.Printf("Failed to get picture stream: %v", err)
		return c.String(404, "Image not found")
	}
	defer reader.Close()

	// Set content type and stream to response
	c.Response().Header().Set("Content-Type", contentType)
	c.Response().Header().Set("Cache-Control", "public, max-age=86400") // Cache for 24 hours
	return c.Stream(200, contentType, reader)
}
