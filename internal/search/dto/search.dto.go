package dto

import (
	"github.com/kien-hoangtrung/github-repository/internal/libs/pagination/dto"
	"github.com/kien-hoangtrung/github-repository/internal/pkg/elasticsearch/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SearchRequestDto struct {
	Keyword string `json:"keyword" validate:"required"`
	Type    string `json:"type"`
	dto.PaginationRequest
}

type ElasticSearchFileResult struct {
	service.ElasticSearchResult
	Hits struct {
		Hits []struct {
			Source struct {
				Id       primitive.ObjectID `json:"id"`
				RepoName string             `json:"repo_name"`
				FileName string             `json:"file_name"`
				Path     string             `json:"path"`
			} `json:"_source"`
			Highlight struct {
				Content []string `json:"content"`
			} `json:"highlight"`
		} `json:"hits"`
	} `json:"hits"`
}

type ElasticSearchRepoResult struct {
	service.ElasticSearchResult
	Hits struct {
		Hits []struct {
			Source struct {
				Id          primitive.ObjectID `json:"id"`
				Name        string             `json:"name"`
				Description string             `json:"description"`
				Topics      []string           `json:"topics"`
			} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

type SearchResponseDto struct {
	dto.PaginationResponse
}
