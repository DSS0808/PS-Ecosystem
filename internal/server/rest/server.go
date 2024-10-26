package rest

import (
	"fmt"
	"net/http"

	_ "github.com/VikaPaz/pantheon/docs"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	router  *gin.Engine
	service Service
}

type Service interface {
}

func NewSrver(svc Service) *Server {
	router := gin.Default()
	return &Server{
		router:  router,
		service: svc,
	}
}

func (h *Server) registerRoutes() {
	tasks := h.router.Group("/app")
	{
		tasks.POST("/", h.main)
	}
	h.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

func (s *Server) Run(port string) {
	s.registerRoutes()
	s.router.Run(":" + port)
}

// func (s *Server) create(c *gin.Context) {
// 	var req CreateRequest

// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		Response(c, nil, http.StatusBadRequest, err)
// 		return
// 	}
// }

// @Summary Send welcome message
// @Description Send welcome message
// @Tags app
// @Accept json
// @Produce json
// @Param request body mainRequest true "msg"
// @Success 200 {object} mainResponse "hello msg"
// @Failure 400
// @Router /app/ [post]
func (s *Server) main(c *gin.Context) {
	var req mainRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		Response(c, nil, http.StatusBadRequest, err)
		return
	}

	res := mainResponse{
		Msg: fmt.Sprintf("Hello %s", req.Msg),
	}

	Response(c, res, http.StatusOK, nil)
}

type mainRequest struct {
	Msg string
}

type mainResponse struct {
	Msg string
}

func Response(
	c *gin.Context,
	responseBody interface{},
	status int,
	err error,
) {
	if err != nil {
		responseBody = gin.H{"error": err.Error()}
	}
	c.JSON(status, responseBody)
}
