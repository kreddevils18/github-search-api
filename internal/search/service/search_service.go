package service

import (
	"github.com/kien-hoangtrung/github-repository/internal/file/entity"
	"github.com/kien-hoangtrung/github-repository/internal/file/service"
	dto2 "github.com/kien-hoangtrung/github-repository/internal/libs/pagination/dto"
	service3 "github.com/kien-hoangtrung/github-repository/internal/pkg/elasticsearch/service"
	"github.com/kien-hoangtrung/github-repository/internal/pkg/log"
	entity2 "github.com/kien-hoangtrung/github-repository/internal/repository/entity"
	service2 "github.com/kien-hoangtrung/github-repository/internal/repository/service"
	"github.com/kien-hoangtrung/github-repository/internal/search/dto"
	"github.com/kien-hoangtrung/github-repository/internal/utils"
	"math"
)

type SearchService struct {
	logger      log.ILogger
	fileService *service.FileService
	repoService *service2.RepositoryService
	esService   *service3.ElasticSearchService
}

func NewSearchService(
	logger log.ILogger,
	fileService *service.FileService,
	repoService *service2.RepositoryService,
	esService *service3.ElasticSearchService,
) *SearchService {
	return &SearchService{logger, fileService, repoService, esService}
}

func (s SearchService) Search(req *dto.SearchRequestDto) (*dto.SearchResponseDto, error) {
	from := (req.Page - 1) * req.Limit

	var err error
	result := new(service3.ElasticSearchResult)
	var data []interface{}
	switch req.Type {
	case "code":
		result, err = s.fileService.Search(from, req.Limit, req.Keyword)
		if err != nil {
			return nil, err
		}
		fileResult := &dto.ElasticSearchFileResult{}
		err = utils.Recast(result, fileResult)
		for _, item := range fileResult.Hits.Hits {
			for _, content := range item.Highlight.Content {
				fileEntity := entity.FileEntity{
					ID:       item.Source.Id,
					RepoName: item.Source.RepoName,
					Path:     item.Source.Path,
					FileName: item.Source.FileName,
					Content:  content,
				}
				data = append(data, fileEntity)
			}
		}
		break
	default:
		result, err = s.repoService.Search(from, req.Limit, req.Keyword)
		if err != nil {
			return nil, err
		}
		repoResult := &dto.ElasticSearchRepoResult{}
		err = utils.Recast(result, repoResult)
		for _, item := range repoResult.Hits.Hits {
			repoEntity := entity2.RepositoryEntity{
				ID:          item.Source.Id,
				Name:        item.Source.Name,
				Description: item.Source.Description,
				Topics:      item.Source.Topics,
			}
			data = append(data, repoEntity)
		}
		break
	}

	totalPages := int64(math.Ceil(float64(result.Hits.Total.Value) / float64(req.Limit)))
	var hasNext bool
	if req.Page < totalPages {
		hasNext = true
	} else {
		hasNext = false
	}

	return &dto.SearchResponseDto{
		PaginationResponse: dto2.PaginationResponse{
			CurrentPage:    req.Page,
			SkippedRecords: req.Skip,
			TotalRecords:   result.Hits.Total.Value,
			TotalPages:     totalPages,
			HasNext:        hasNext,
			Data:           data,
		},
	}, nil
}
