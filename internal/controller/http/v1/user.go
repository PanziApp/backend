package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/PanziApp/backend/internal/domain"
	"github.com/PanziApp/backend/internal/usecase"
	"github.com/PanziApp/backend/pkg/logger"
)

type userRoutes struct {
	t usecase.Translation
	l logger.Interface
}

func newUserRoutes(handler *gin.RouterGroup, t usecase.Translation, l logger.Interface) {
	r := &userRoutes{t, l}

	handler.POST("/sign-up")
	handler.POST("/sign-in")
	handler.POST("/reset-password/link")
	handler.POST("/reset-password")

	h := handler.Group("/users")
	{
		h.POST("/sign-out")
		h.GET("/profile")
		h.POST("/profile")
		h.POST("/password")
		h.GET("/avatar")
		h.POST("/avatar")
	}
}

type historyResponse struct {
	History []domain.Translation `json:"history"`
}

// @Summary     Show history
// @Description Show all translation history
// @ID          history
// @Tags  	    translation
// @Accept      json
// @Produce     json
// @Success     200 {object} historyResponse
// @Failure     500 {object} response
// @Router      /translation/history [get]
func (r *userRoutes) history(c *gin.Context) {
	translations, err := r.t.History(c.Request.Context())
	if err != nil {
		r.l.Error(err, "http - v1 - history")
		errorResponse(c, http.StatusInternalServerError, "database problems")

		return
	}

	c.JSON(http.StatusOK, historyResponse{translations})
}

type doTranslateRequest struct {
	Source      string `json:"source"       binding:"required"  example:"auto"`
	Destination string `json:"destination"  binding:"required"  example:"en"`
	Original    string `json:"original"     binding:"required"  example:"текст для перевода"`
}

// @Summary     Translate
// @Description Translate a text
// @ID          do-translate
// @Tags  	    translation
// @Accept      json
// @Produce     json
// @Param       request body doTranslateRequest true "Set up translation"
// @Success     200 {object} domain.Translation
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /translation/do-translate [post]
func (r *userRoutes) doTranslate(c *gin.Context) {
	var request doTranslateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - v1 - doTranslate")
		errorResponse(c, http.StatusBadRequest, "invalid request body")

		return
	}

	translation, err := r.t.Translate(
		c.Request.Context(),
		domain.Translation{
			Source:      request.Source,
			Destination: request.Destination,
			Original:    request.Original,
		},
	)
	if err != nil {
		r.l.Error(err, "http - v1 - doTranslate")
		errorResponse(c, http.StatusInternalServerError, "translation service problems")

		return
	}

	c.JSON(http.StatusOK, translation)
}
