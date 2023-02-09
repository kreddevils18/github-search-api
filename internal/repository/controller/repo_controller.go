package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kien-hoangtrung/github-repository/internal/libs/pagination"
	"github.com/kien-hoangtrung/github-repository/internal/libs/pagination/dto"
	"github.com/kien-hoangtrung/github-repository/internal/pkg/log"
	"github.com/kien-hoangtrung/github-repository/internal/pkg/validation"
	"github.com/kien-hoangtrung/github-repository/internal/repository/service"
)

type RepoController struct {
	logger      log.ILogger
	repoService *service.RepositoryService
}

func NewRepoController(logger log.ILogger, repoService *service.RepositoryService) *RepoController {
	return &RepoController{logger, repoService}
}

func (r *RepoController) List(ctx *fiber.Ctx) error {
	r.logger.Info("RepoController: List")
	listRequest := new(dto.PaginationRequest)

	if err := ctx.QueryParser(listRequest); err != nil {
		r.logger.Infof("RepoController: ErrorBodyParseList")
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	errors := validation.ValidateStruct(listRequest)
	if errors != nil {
		r.logger.Infof("RepoController: ErrorValidationStruct")
		return ctx.Status(fiber.StatusBadRequest).JSON(errors)
	}

	res, err := r.repoService.List(pagination.PaginationParamsTransform(listRequest))
	if err != nil {
		r.logger.Infof("RepoController: ErrorRepoServiceList")
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(res)
}
