package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/kien-hoangtrung/github-repository/internal/file/repository"
	service2 "github.com/kien-hoangtrung/github-repository/internal/file/service"
	"github.com/kien-hoangtrung/github-repository/internal/pkg/config"
	"github.com/kien-hoangtrung/github-repository/internal/pkg/dbcontext"
	service3 "github.com/kien-hoangtrung/github-repository/internal/pkg/elasticsearch/service"
	"github.com/kien-hoangtrung/github-repository/internal/pkg/github"
	"github.com/kien-hoangtrung/github-repository/internal/pkg/log"
	repository2 "github.com/kien-hoangtrung/github-repository/internal/repository/repository"
	"github.com/kien-hoangtrung/github-repository/internal/repository/service"
	"os"
	"strconv"
	"sync"
)

func main() {
	err := godotenv.Load()

	conf, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	logger, err := log.NewLogger(conf)
	if err != nil {
		panic(err)
	}

	logger.Info("seeding...")

	// DBContext
	db, err := dbcontext.NewMongo(conf)
	if err != nil {
		logger.Fatal(err)
	}

	// Repository
	fileRepository := repository.NewFileRepository(logger, db)
	repoRepository := repository2.NewRepoRepository(logger, db)

	// Service
	githubService := github.NewGithub(conf, logger)
	elasticSearchService, err := service3.NewElasticSearchService(conf, logger)
	if err != nil {
		logger.Fatal(err)
	}
	repoService := service.NewRepositoryService(logger, repoRepository, elasticSearchService)
	fileService := service2.NewFileService(logger, fileRepository, elasticSearchService)

	// Execute
	var wg sync.WaitGroup
	page, _ := strconv.Atoi(os.Getenv("PAGE"))
	ctx := context.Background()

	repos, err := githubService.ListRepos(ctx, conf.Organization.Team, page)
	if err != nil {
		logger.Fatal(err)
	}

	for len(repos) > 0 {
		wg.Add(1)
		go func() {
			for _, repo := range repos {
				err = repoService.SaveRepository(repo)
				if err != nil {
					logger.Info(err)
				}
				err = fileService.SaveRepositoryContent(repo)
				if err != nil {
					logger.Info(err)
				}
			}

			wg.Done()
		}()

		page += 1
		repos, err = githubService.ListRepos(ctx, conf.Organization.Team, page)
		if err != nil {
			logger.Fatal(err)
		}
	}

	wg.Wait()
}
