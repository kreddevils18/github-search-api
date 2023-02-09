package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/go-github/v49/github"
	"github.com/kien-hoangtrung/github-repository/internal/file/entity"
	"github.com/kien-hoangtrung/github-repository/internal/file/repository"
	"github.com/kien-hoangtrung/github-repository/internal/pkg/config"
	"github.com/kien-hoangtrung/github-repository/internal/pkg/elasticsearch/service"
	"github.com/kien-hoangtrung/github-repository/internal/pkg/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type FileService struct {
	logger        log.ILogger
	repository    *repository.FileRepository
	elasticsearch *service.ElasticSearchService
}

func NewFileService(
	logger log.ILogger,
	repository *repository.FileRepository,
	elasticsearch *service.ElasticSearchService,
) *FileService {
	return &FileService{
		logger,
		repository,
		elasticsearch}
}

func (f FileService) SaveRepositoriesContent(repos []*github.Repository) error {
	for len(repos) > 0 {
		for _, repo := range repos {
			err := f.SaveRepositoryContent(repo)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (f FileService) SaveRepositoryContent(repo *github.Repository) error {
	f.logger.Info(repo.GetName())
	tmpPath := config.GetSourcePath() + "/../../../tmp///"
	filePath := fmt.Sprintf("%s/%s", filepath.Dir(tmpPath), repo.GetName())
	cmd := exec.Command("git", "clone", repo.GetSSHURL(), filePath)
	err := cmd.Run()

	if err != nil {
		f.logger.Infof("error cmd command %s : %s", cmd.String(), err)
	}

	err = filepath.Walk(filePath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			f.logger.Infof("walk err %s", err)
			return err
		}
		var extensions = map[string]bool{
			".mp3":     true,
			".mp4":     true,
			".png":     true,
			".jpg":     true,
			".jpeg":    true,
			".otf":     true,
			".woff":    true,
			".ttf":     true,
			".gif":     true,
			".bmglyph": true,
			".webm":    true,
			".ktx":     true,
			".ico":     true,
		}

		if !info.IsDir() && info.Size() > 0 && !strings.Contains(path, ".git") && !extensions[filepath.Ext(path)] {
			file, err := os.ReadFile(path)
			if err != nil {
				f.logger.Infof("error opening %s: %s", path, err)
			}

			// Filepath transform
			re := regexp.MustCompile(fmt.Sprintf(".*/%s/(.*)", repo.GetName()))
			match := re.FindStringSubmatch(path)
			tmp := strings.Split(match[1], "/")
			fileName := tmp[len(tmp)-1]
			fileHttpPath := fmt.Sprintf("%s%s%s", repo.GetHTMLURL(), "/blob/master/", match[1])
			f.logger.Info(fileHttpPath)

			createdFile, err := f.Create(*repo.Name, fileName, fileHttpPath, string(file))
			if err != nil {
				f.logger.Infof("error create file %s : %d %s", path, info.Size(), err)
			} else {
				err = f.elasticsearch.IndexRequest("file", createdFile.ID.Hex(), *createdFile)
				if err != nil {
					f.logger.Infof("error elastic index %s: %s", path, err)
				}
			}

		}
		return nil
	})
	if err != nil {
		return err
	}

	err = os.RemoveAll(filePath)
	if err != nil {
		f.logger.Fatal(err)
	}

	return nil
}

func (f FileService) Create(repoName string, fileName string, path string, content string) (*entity.FileEntity, error) {
	_, err := f.repository.FindByPath(path)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return nil, err
		}
	}

	if err == mongo.ErrNoDocuments {
		file := entity.FileEntity{
			ID:       primitive.NewObjectID(),
			RepoName: repoName,
			FileName: fileName,
			Path:     path,
			Content:  content,
		}

		err = f.repository.CreateOne(file)
		if err != nil {
			return nil, err
		}
		return &file, nil
	} else {
		return nil, errors.New("file is existed")
	}
}

func (f FileService) Search(from int64, size int64, keyword string) (*service.ElasticSearchResult, error) {
	query := fmt.Sprintf(
		`{
      "_source": ["id", "repo_name", "file_name", "path"], 
      "from": %d,
      "size": %d,
      "query": {
        "match_phrase": {
          "content": "%s"
        }
      },
      "highlight": {
        "pre_tags": ["<span>"],
        "post_tags": ["</span>"],
        "fields": {
          "content": {}
        }
      }
    }`, from, size, keyword)

	res, err := f.elasticsearch.Query("file", query)
	if err != nil {
		f.logger.Infof("FileService: ErrorSearch: %s", err)
		return nil, err
	}
	defer res.Body.Close()

	result := new(service.ElasticSearchResult)
	if err := json.NewDecoder(res.Body).Decode(result); err != nil {
		f.logger.Infof("FileService: ErrorElasticsearchBodyDecode: %s", err)
		return nil, err
	}

	return result, nil
}
