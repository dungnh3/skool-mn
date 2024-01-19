package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/dungnh3/skool-mn/config"
	"github.com/dungnh3/skool-mn/docs"
	"github.com/dungnh3/skool-mn/internal/middleware"
	"github.com/dungnh3/skool-mn/internal/repositories"
	l "github.com/dungnh3/skool-mn/pkg/log"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"net/http"
)

type Server struct {
	server *http.Server
	r      repositories.Repository
	cfg    *config.Config
	logger l.Logger
}

func New(cfg *config.Config, r repositories.Repository) *Server {
	logger := l.New().Named("server")
	port := cfg.Server.HTTP.Port
	router := gin.Default()
	router.RedirectTrailingSlash = true
	router.RedirectFixedPath = true
	router.HandleMethodNotAllowed = true
	router.RemoveExtraSlash = true

	s := &Server{
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: router,
		},
		r:      r,
		cfg:    cfg,
		logger: logger,
	}
	p := ginprometheus.NewPrometheus("gin")
	p.Use(router)

	router.GET("/health", s.healthCheck)
	router.GET("/live", s.liveCheck)

	docs.SwaggerInfo.BasePath = ""
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	router.Use(middleware.Auth)
	parentRoute := router.Group("/parents")
	{
		parentRoute.GET("/:id/students", s.listStudents)
		parentRoute.POST("/register", s.registerPickUpTime)
		parentRoute.PUT("/registers/:id/waiting", s.waitingFromParent)
		parentRoute.PUT("/registers/:id/confirm", s.confirmCompleted)
		parentRoute.PUT("/registers/:id/cancel", s.cancelFromParent)
	}

	teacherRoute := router.Group("/teachers")
	{
		teacherRoute.PUT("/registers/:id/confirm", s.confirmFromTeacher)
		teacherRoute.PUT("/registers/:id/reject", s.rejectFromTeacher)
	}

	studentsRoute := router.Group("/students")
	{
		studentsRoute.PUT("/:id/leave", s.studentLeaveClass)
		studentsRoute.PUT("/:id/out", s.studentOutSchool)
	}
	return s
}

// Run application
func (s *Server) Run() error {
	s.logger.Info("Start the server at", l.Object("address", s.server.Addr))
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

// Close app and all the resources
func (s *Server) Close(ctx context.Context) error {
	s.logger.Info("Shutting down server...")
	return s.server.Shutdown(ctx)
}

func (s *Server) healthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "OK")
}

func (s *Server) liveCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "OK")
}
