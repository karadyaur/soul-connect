package routers

import (
	"github.com/gin-gonic/gin"
	"soul-connect/sc-api-getaway/internal/controllers"
)

type postRouter struct {
	controller *controllers.PostController
}

func newPostRouter(controller *controllers.PostController) *postRouter {
	return &postRouter{controller: controller}
}

func (r *postRouter) setPostRoutes(group *gin.RouterGroup) {
	group.POST("/posts", r.controller.CreatePost)
	group.GET("/posts", r.controller.ListPosts)
	group.GET("/posts/:post_id", r.controller.GetPost)
	group.PUT("/posts/:post_id", r.controller.UpdatePost)
	group.DELETE("/posts/:post_id", r.controller.DeletePost)

	group.POST("/posts/:post_id/comments", r.controller.AddComment)
	group.GET("/posts/:post_id/comments", r.controller.ListComments)

	group.POST("/posts/:post_id/likes", r.controller.LikePost)
	group.DELETE("/posts/:post_id/likes", r.controller.UnlikePost)

	group.POST("/comments/:comment_id/likes", r.controller.LikeComment)
	group.DELETE("/comments/:comment_id/likes", r.controller.UnlikeComment)

	group.GET("/labels", r.controller.ListLabels)
	group.POST("/posts/:post_id/labels", r.controller.AddLabelToPost)
	group.DELETE("/posts/:post_id/labels/:label_id", r.controller.RemoveLabelFromPost)
}
