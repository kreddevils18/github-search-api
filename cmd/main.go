package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/kien-hoangtrung/github-repository/internal/file/repository"
	service2 "github.com/kien-hoangtrung/github-repository/internal/file/service"
	"github.com/kien-hoangtrung/github-repository/internal/pkg/config"
	"github.com/kien-hoangtrung/github-repository/internal/pkg/dbcontext"
	"github.com/kien-hoangtrung/github-repository/internal/pkg/elasticsearch/service"
	"github.com/kien-hoangtrung/github-repository/internal/pkg/log"
	controller2 "github.com/kien-hoangtrung/github-repository/internal/repository/controller"
	repository2 "github.com/kien-hoangtrung/github-repository/internal/repository/repository"
	service3 "github.com/kien-hoangtrung/github-repository/internal/repository/service"
	"github.com/kien-hoangtrung/github-repository/internal/search/controller"
	service4 "github.com/kien-hoangtrung/github-repository/internal/search/service"
)

func main() {
	// Initial
	conf, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	logger, err := log.NewLogger(conf)
	if err != nil {
		panic(err)
	}

	// Repositories
	db, err := dbcontext.NewMongo(conf)
	fileRepository := repository.NewFileRepository(logger, db)
	repoRepository := repository2.NewRepoRepository(logger, db)

	// Services
	elasticsearch, err := service.NewElasticSearchService(conf, logger)
	if err != nil {
		logger.Fatal(err)
	}
	fileService := service2.NewFileService(logger, fileRepository, elasticsearch)
	repoService := service3.NewRepositoryService(logger, repoRepository, elasticsearch)
	searchService := service4.NewSearchService(logger, fileService, repoService, elasticsearch)

	// Controllers
	//fileController := controller.NewFileController(logger, fileService)
	repoController := controller2.NewRepoController(logger, repoService)
	searchController := controller.NewSearchController(logger, searchService)

	// Routes define
	app := fiber.New()
	api := app.Group("api")

	// Middlewares
	app.Use(cors.New())

	//fileApi := api.Group("files")
	//fileApi.Get("/search", fileController.Search)
	api.Get("/search", searchController.Search)

	repoApi := api.Group("repos")
	repoApi.Get("/", repoController.List)

	logger.Fatal(app.Listen(conf.Server.Port))
}
