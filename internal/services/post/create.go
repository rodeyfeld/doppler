package post

import "github.com/rodeyfeld/doppler/internal/models"

func (postService *Service) Create(post *models.Post) {
	postService.DB.Create(post)
}
