package controller

import (
	"github.com/kien-hoangtrung/github-repository/internal/file/service"
	"github.com/kien-hoangtrung/github-repository/internal/pkg/log"
)

type FileController struct {
	logger      log.ILogger
	fileService *service.FileService
}

func NewFileController(logger log.ILogger, fileService *service.FileService) *FileController {
	return &FileController{logger, fileService}
}
