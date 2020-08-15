package http

import (
	atDomain "github.com/carloshjoaquim/bookstore-oauth-api/src/domain/access_token"
	"github.com/carloshjoaquim/bookstore-oauth-api/src/services/access_token"
	"github.com/carloshjoaquim/bookstore-oauth-api/src/utils/errors_utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type AccessTokenHandler interface {
	GetById(*gin.Context)
	Create(ctx *gin.Context)
}

type accessTokenHandler struct {
	service access_token.Service
}

func NewHandler(service access_token.Service) AccessTokenHandler {
	return &accessTokenHandler{
		service: service,
	}
}

func (handler *accessTokenHandler) GetById(c *gin.Context) {
	accessTokenId := strings.TrimSpace(c.Param("access_token_id"))

	accessToken, err := handler.service.GetById(accessTokenId)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusOK, accessToken)
}

func (handler *accessTokenHandler) Create(c *gin.Context) {
	var request atDomain.AccessTokenRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		restErr := errors_utils.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}

	accessToken, err := handler.service.Create(request)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusCreated, accessToken)
}
