package http

import (
	"net/http"

	"github.com/filhodanuvem/ytgoapi/internal"
	"github.com/filhodanuvem/ytgoapi/internal/database"
	"github.com/filhodanuvem/ytgoapi/internal/post"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	response, err := service.Create(ctx, post)

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
	id, err := uuid.Parse(param)
	if err != nil {
		ctx.JSON(http.StatusNotFound, nil)
	}

	if err := service.Delete(ctx, id); err != nil {
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
	id, err := uuid.Parse(param)
	if err != nil {
		ctx.JSON(http.StatusNotFound, nil)
	}

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
