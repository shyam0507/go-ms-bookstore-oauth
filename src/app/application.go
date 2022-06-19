package app

import (
	"github.com/gin-gonic/gin"
	"github.com/shyam0507/go-ms-bookstore-oauth/src/http"
	"github.com/shyam0507/go-ms-bookstore-oauth/src/repository/db"
	"github.com/shyam0507/go-ms-bookstore-oauth/src/repository/rest"
	"github.com/shyam0507/go-ms-bookstore-oauth/src/services/access_token"
)

var router = gin.Default()

func StartApplication() {

	atService := access_token.NewService(rest.NewRepository(), db.New())
	atHandler := http.NewHandler(atService)

	router.GET("/oauth/access_token/:access_token_id", atHandler.GetById)
	router.POST("/oauth/access_token", atHandler.Create)

	router.Run(":8080")
}
