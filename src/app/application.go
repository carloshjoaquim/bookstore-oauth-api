package app

import (
	"github.com/carloshjoaquim/bookstore-oauth-api/src/http"
	"github.com/carloshjoaquim/bookstore-oauth-api/src/repository/db"
	"github.com/carloshjoaquim/bookstore-oauth-api/src/repository/rest"
	"github.com/carloshjoaquim/bookstore-oauth-api/src/services/access_token"
	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

func StartApplication() {
	atHandler := http.NewHandler(
		access_token.NewService(rest.NewUsersRepository(), db.NewRepository()))

	router.GET("/oauth/access_token/:access_token_id", atHandler.GetById)
	router.POST("/oauth/access_token/", atHandler.Create)

	router.Run(":8080")
}
