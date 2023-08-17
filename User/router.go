package user

import "github.com/gin-gonic/gin"

type Controller struct {
	CreateStrategy UserHandler
	DeleteStrategy UserHandler
	GetAllStrategy UserHandler
	GetOneStrategy UserHandler
	UpdateStrategy UserHandler
}

func NewUserController(
	createStrategy UserHandler,
	deleteStrategy UserHandler,
	getAllStrategy UserHandler,
	getOneStrategy UserHandler,
	updateStrategy UserHandler,
) *Controller {
	return &Controller{
		CreateStrategy: createStrategy,
		DeleteStrategy: deleteStrategy,
		GetAllStrategy: getAllStrategy,
		GetOneStrategy: getOneStrategy,
		UpdateStrategy: updateStrategy,
	}
}

func (ctrl *Controller) CreateUserHandler(c *gin.Context) {
	ctrl.CreateStrategy.Handle(c)
}
func (ctrl *Controller) DeleteUserHandler(c *gin.Context) {
	ctrl.DeleteStrategy.Handle(c)
}
func (ctrl *Controller) GetAllUserHandler(c *gin.Context) {
	ctrl.GetAllStrategy.Handle(c)
}
func (ctrl *Controller) GetOneUserHandler(c *gin.Context) {
	ctrl.GetOneStrategy.Handle(c)
}
func (ctrl *Controller) UpdateUserHandler(c *gin.Context) {
	ctrl.UpdateStrategy.Handle(c)
}

func SetupUserController(repo UserRepository) *Controller {
	createStrategy := &CreateUserStrategy{Repo: repo}
	deleteStrategy := &DeleteUserStrategy{Repo: repo}
	getAllStrategy := &GetAllUserStrategy{Repo: repo}
	getOneStrategy := &GetSpesificUserStrategy{Repo: repo}
	updateStrategy := &UpdateUserStrategy{Repo: repo}

	controller := NewUserController(
		createStrategy,
		deleteStrategy,
		getAllStrategy,
		getOneStrategy,
		updateStrategy,
	)

	return controller
}

func SetupUserRoutes(r *gin.Engine, userController *Controller) {

	users := r.Group("/users")
	{
		users.POST("/", userController.CreateUserHandler)
		users.GET("/:id", userController.GetOneUserHandler)
		users.GET("/", userController.GetAllUserHandler)
		users.PUT("/:id", userController.UpdateUserHandler)
		users.DELETE("/:id", userController.DeleteUserHandler)
	}
}
