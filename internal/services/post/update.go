package post

import (
	"github.com/rodeyfeld/doppler/internal/models"
	"github.com/rodeyfeld/doppler/internal/requests"
)

func (postService *Service) Update(post *models.Post, updatePostRequest *requests.UpdatePostRequest) {
	post.Content = updatePostRequest.Content
	post.Title = updatePostRequest.Title
	postService.DB.Save(post)
}
