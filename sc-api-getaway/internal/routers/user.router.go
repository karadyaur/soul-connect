package routers

import (
	"github.com/gin-gonic/gin"
	"soul-connect/sc-api-getaway/internal/config"
	"soul-connect/sc-api-getaway/internal/controllers"
)

type userRouter struct {
	controller *controllers.UserController
	config     *config.Config
}

func newUserRouter(controller *controllers.UserController, config *config.Config) *userRouter {
	return &userRouter{controller: controller, config: config}
}

func (ur *userRouter) setUserRoutes(rg *gin.RouterGroup) {
	router := rg.Group("users")
	router.POST("", ur.controller.CreateProfile)
	router.GET("/:id", ur.controller.GetProfile)
	router.PUT("/:id", ur.controller.UpdateProfile)
	router.DELETE("/:id", ur.controller.DeleteProfile)
	router.POST("/:id/subscriptions", ur.controller.Subscribe)
	router.GET("/:id/subscriptions", ur.controller.ListSubscriptions)
	router.DELETE("/:id/subscriptions/:author_id", ur.controller.Unsubscribe)
}
