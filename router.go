package main

import (
	"EmptyClassroom/logs"
	"embed"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"io/fs"
	"net/http"
)

//go:embed frontend/dist
var f embed.FS

type embedFileSystem struct {
	http.FileSystem
}

func (e embedFileSystem) Exists(prefix string, path string) bool {
	_, err := e.Open(path)
	if err != nil {
		return false
	}
	return true
}

func EmbedFolder(fsEmbed embed.FS, targetPath string) static.ServeFileSystem {
	fsys, err := fs.Sub(fsEmbed, targetPath)
	if err != nil {
		panic(err)
	}
	return embedFileSystem{
		FileSystem: http.FS(fsys),
	}
}

func SetRouter(r *gin.Engine) {
	r.Use(static.Serve("/", EmbedFolder(f, "frontend/dist")))

	apiGroup := r.Group("/api").Use(logs.SetNewContextForGinContext)
	{
		apiGroup.GET("/get_data", GetData)
		apiGroup.POST("/report", Report)
	}
}
