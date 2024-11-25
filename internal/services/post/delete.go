package post

import "github.com/rodeyfeld/doppler/internal/models"

func (postService *Service) Delete(post *models.Post) {
	postService.DB.Delete(post)
}
