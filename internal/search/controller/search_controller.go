package controller

import (
  "github.com/gofiber/fiber/v2"
  "github.com/kien-hoangtrung/github-repository/internal/libs/pagination"
  "github.com/kien-hoangtrung/github-repository/internal/pkg/log"
  "github.com/kien-hoangtrung/github-repository/internal/pkg/validation"
  "github.com/kien-hoangtrung/github-repository/internal/search/dto"
  "github.com/kien-hoangtrung/github-repository/internal/search/service"
)

type SearchController struct {
  logger        log.ILogger
  searchService *service.SearchService
}

func NewSearchController(logger log.ILogger, searchService *service.SearchService) *SearchController {
  return &SearchController{logger, searchService}
}

func (s *SearchController) Search(ctx *fiber.Ctx) error {
  requestDto := new(dto.SearchRequestDto)

  if err := ctx.QueryParser(requestDto); err != nil {
    return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
      "message": err.Error(),
    })
  }

  s.logger.Infof("Search with params: %v", *requestDto)
  errors := validation.ValidateStruct(*requestDto)
  if errors != nil {
    return ctx.Status(fiber.StatusBadRequest).JSON(errors)
  }

  res, err := s.searchService.Search(pagination.PaginationParamsTransform(requestDto))
  if err != nil {
    return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
      "message": err.Error(),
    })
  }

  return ctx.JSON(res)
}
