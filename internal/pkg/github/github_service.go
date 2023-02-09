package github

import (
	"context"
	"github.com/google/go-github/v49/github"
	"github.com/kien-hoangtrung/github-repository/internal/pkg/config"
	"github.com/kien-hoangtrung/github-repository/internal/pkg/log"
	"golang.org/x/oauth2"
)

type GithubService struct {
	logger *log.Logger
	config *config.Config
	client *github.Client
}

func NewGithub(conf *config.Config, logger *log.Logger) *GithubService {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: conf.Github.ApiKey})
	tc := oauth2.NewClient(context.TODO(), ts)
	client := github.NewClient(tc)
	return &GithubService{logger, conf, client}
}

func (g *GithubService) ListRepos(ctx context.Context, slug string, page int) ([]*github.Repository, error) {
	options := &github.ListOptions{
		Page:    page,
		PerPage: 5,
	}
	repos, _, err := g.client.Teams.ListTeamReposBySlug(ctx, g.config.Organization.Name, slug, options)
	if err != nil {
		return nil, err
	}

	return repos, nil
}

func (g GithubService) GetTree(ctx context.Context, repoName string) (*github.Tree, error) {
	tree, _, err := g.client.Git.GetTree(ctx, g.config.Organization.Name, repoName, "master", true)
	if err != nil {
		return nil, err
	}

	return tree, nil
}

func (g GithubService) GetContents(ctx context.Context, repoName string, path string) (*github.RepositoryContent, error) {
	fileContent, _, _, err := g.client.Repositories.GetContents(ctx, g.config.Organization.Name, repoName, path, nil)
	if err != nil {
		return nil, err
	}

	return fileContent, nil
}
