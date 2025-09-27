package routes

import (
	"frog-go/internal/core/ports/outbound/messagebus"
	"frog-go/internal/core/ports/outbound/repository"
	"frog-go/internal/core/service"

	"frog-go/internal/core/service/upload"
	"frog-go/internal/http/handler"
	"frog-go/internal/http/middlewares"
	"frog-go/internal/utils/logger"

	"github.com/gin-gonic/gin"
)

type Router struct {
	log  *logger.Logger
	repo repository.Repository
	mbus messagebus.MessageBus
}

func NewRouter(log *logger.Logger, repo repository.Repository, mbus messagebus.MessageBus) *Router {
	return &Router{
		log:  log,
		repo: repo,
		mbus: mbus,
	}
}

func (r *Router) Setup(enableDebug bool) *gin.Engine {
	if enableDebug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.Default()

	v1 := engine.Group("/api/v1")

	v1.Use(middlewares.ErrorMiddleware(r.log))
	// v1.Use(middlewares.CORSMiddleware())
	v1.Use(middlewares.UUIDMiddleware(r.log))

	authService := service.NewAuthService(r.repo)
	userService := service.NewUserService(r.repo)
	authHandler := handler.NewAuthHandler(authService, userService)

	registerAuthRoutes(engine.Group("/api/auth"), authHandler)

	v1.Use(middlewares.AuthMiddleware(r.log, authService.ValidateToken))

	engine.StaticFile("/favicon.ico", "./static/favicon.ico")
	registerDocsRoutes(engine.Group("/docs/v1"))

	transactionService := service.NewTransactionService(r.repo)
	transactionHandler := handler.NewTransactionHandler(transactionService)
	registerTransactionRoutes(v1.Group("/transactions"), transactionHandler)

	invoiceService := service.NewInvoiceService(r.repo)
	invoiceHandler := handler.NewInvoiceHandler(invoiceService)
	registerInvoiceRoutes(v1.Group("/invoices"), invoiceHandler)

	categoryService := service.NewCategoryService(r.repo)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	registerCategoryRoutes(v1.Group("/categories"), categoryHandler)

	uploadService := upload.NewUploadService(r.mbus)
	uploadHander := handler.NewUploadHandler(uploadService)
	registerUploadRoutes(v1.Group("/upload"), uploadHander)

	return engine
}

func registerDocsRoutes(router *gin.RouterGroup) {
	router.StaticFile("/swagger.json", "./docs/v1/swagger.json")

	router.GET("", func(c *gin.Context) {
		c.File("./docs/v1/docs.html")
	})
	router.GET("/redoc", func(c *gin.Context) {
		c.File("./docs/v1/redoc.html")
	})
	router.GET("/scalar", func(c *gin.Context) {
		c.File("./docs/v1/scalar.html")
	})
	router.GET("/swagger", func(c *gin.Context) {
		c.File("./docs/v1/swagger.html")
	})

}

func registerAuthRoutes(router *gin.RouterGroup, handler *handler.AuthHandler) {
	router.POST("/login", handler.Login)
}

func registerTransactionRoutes(router *gin.RouterGroup, handler *handler.TransactionHandler) {
	router.POST("", handler.CreateTransactionHandler)
	router.GET("", handler.ListTransactionsHandler)
	router.GET("/:id", handler.GetTransactionByIDHandler)
	router.PUT("/:id", handler.UpdateTransactionHandler)
	router.DELETE("/:id", handler.DeleteTransactionHandler)
	router.GET("/summary", handler.TransactionsSummaryHandler)
	router.GET("/stats", handler.TransactionsGeneralStatsHandler)
}

func registerInvoiceRoutes(router *gin.RouterGroup, handler *handler.InvoiceHandler) {
	router.POST("", handler.CreateInvoiceHandler)
	router.GET("", handler.ListInvoicesHandler)
	router.GET("/:id", handler.GetInvoiceByIDHandler)
	router.PUT("/:id", handler.UpdateInvoiceHandler)
	router.DELETE("/:id", handler.DeleteInvoiceHandler)
	router.GET("/:id/transactions", handler.ListInvoiceDebtsHandler)
}

func registerCategoryRoutes(router *gin.RouterGroup, handler *handler.CategoryHandler) {
	router.POST("", handler.CreateCategoryHandler)
	router.GET("", handler.ListCategorysHandler)
	router.GET("/:id", handler.GetCategoryByIDHandler)
	router.PUT("/:id", handler.UpdateCategoryHandler)
	router.DELETE("/:id", handler.DeleteCategoryHandler)
}

func registerUploadRoutes(router *gin.RouterGroup, handler *handler.UploadHandler) {
	router.POST("", handler.ProcessFileHandler)
}
