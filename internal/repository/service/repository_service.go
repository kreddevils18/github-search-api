package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/go-github/v49/github"
	"github.com/kien-hoangtrung/github-repository/internal/libs/pagination/dto"
	"github.com/kien-hoangtrung/github-repository/internal/pkg/elasticsearch/service"
	"github.com/kien-hoangtrung/github-repository/internal/pkg/log"
	"github.com/kien-hoangtrung/github-repository/internal/repository/entity"
	"github.com/kien-hoangtrung/github-repository/internal/repository/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"math"
)

type RepositoryService struct {
	logger         log.ILogger
	repoRepository *repository.RepoRepository
	elasticsearch  *service.ElasticSearchService
}

func NewRepositoryService(
	logger log.ILogger,
	repoRepository *repository.RepoRepository,
	elasticsearch *service.ElasticSearchService) *RepositoryService {
	return &RepositoryService{logger, repoRepository, elasticsearch}
}

func (r RepositoryService) SaveRepositories(repos []*github.Repository) error {
	for _, repo := range repos {
		err := r.SaveRepository(repo)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r RepositoryService) SaveRepository(repo *github.Repository) error {
	createdRepo, err := r.Create(repo.GetName(), repo.GetDescription(), repo.Topics)
	if err != nil {
		r.logger.Info(err)
	} else {
		err = r.elasticsearch.IndexRequest("repository", createdRepo.ID.Hex(), *createdRepo)
		if err != nil {
			r.logger.Infof("error elastic index %s: %s", createdRepo.Name, err)
		}
	}

	return nil
}

func (r RepositoryService) Create(name string, description string, topics []string) (*entity.RepositoryEntity, error) {
	_, err := r.repoRepository.FindByName(name)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return nil, err
		}
	}

	if err == mongo.ErrNoDocuments {
		repoEntity := entity.RepositoryEntity{
			ID:          primitive.NewObjectID(),
			Description: description,
			Name:        name,
			Topics:      topics,
		}

		err = r.repoRepository.CreateOne(repoEntity)
		if err != nil {
			return nil, err
		}
		return &repoEntity, nil
	} else {
		return nil, errors.New("repository is existed")
	}
}

func (r RepositoryService) List(req *dto.PaginationRequest) (*dto.PaginationResponse, error) {
	repoEntities, totalRecords, err := r.repoRepository.FindMany(req.Skip, req.Limit)
	if err != nil {
		return nil, err
	}

	totalPages := int64(math.Ceil(float64(totalRecords) / float64(req.Limit)))

	var hasNext bool
	if req.Page < totalPages {
		hasNext = true
	} else {
		hasNext = false
	}

	res := dto.PaginationResponse{
		CurrentPage:    req.Page,
		SkippedRecords: req.Skip,
		TotalRecords:   totalRecords,
		TotalPages:     totalPages,
		HasNext:        hasNext,
		Data:           repoEntities,
	}

	return &res, nil
}

func (r RepositoryService) Search(from int64, size int64, keyword string) (*service.ElasticSearchResult, error) {
	query := fmt.Sprintf(
		`{
      "from": %d,
      "size": %d,
      "query": {
        "match_phrase": {
          "name": {
            "query": "%s"
          }
        }
      }
    }`, from, size, keyword)

	res, err := r.elasticsearch.Query("repository", query)
	if err != nil {
		r.logger.Infof("RepositoryService: ErrorSearch: %s", err)
		return nil, err
	}
	defer res.Body.Close()

	result := new(service.ElasticSearchResult)
	if err := json.NewDecoder(res.Body).Decode(result); err != nil {
		r.logger.Infof("RepositoryService: ErrorElasticsearchBodyDecode: %s", err)
		return nil, err
	}

	return result, nil
}
