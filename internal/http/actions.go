package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/filhodanuvem/ytgoapi/internal"
	"github.com/filhodanuvem/ytgoapi/internal/database"
	"github.com/filhodanuvem/ytgoapi/internal/post"
	"github.com/gin-gonic/gin"
)

var service post.Service

func Configure() {
	service = post.Service{
		Repository: &post.RepositoryPostgres{
			Conn: database.Conn,
		},
	}
}

func PostPosts(ctx *gin.Context) {
	var post internal.Post
	if err := ctx.BindJSON(&post); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	response, err := service.Create(ctxTimeout, post)


	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func DeletePosts(ctx *gin.Context) {
	param := ctx.Param("id")

	if param == "" {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": post.ErrIdEmpty,
		})
		return
	}


	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := service.Delete(ctxTimeout, param); err != nil {

		statusCode := http.StatusInternalServerError
		if err == post.ErrPostNotFound {
			statusCode = http.StatusNotFound
		}

		ctx.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func GetPosts(ctx *gin.Context) {
	param := ctx.Param("id")

	if param == "" {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": post.ErrIdEmpty,
		})
		return
	}


	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	p, err := service.FindOneByID(ctxTimeout, param)

	p, err := service.FindOneByID(ctx, id)

	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == post.ErrPostNotFound {
			statusCode = http.StatusNotFound
		}

		ctx.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, p)
}

func GetAll(ctx *gin.Context) {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	posts, err := service.FindAll(ctxTimeout)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": http.StatusText(http.StatusInternalServerError),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"posts": posts})
}

func Update(ctx *gin.Context) {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var params internal.ParamsUpdatePost

	if err := ctx.BindJSON(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := service.Update(ctxTimeout, &params); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("user %s updated", params.ID)})
}
