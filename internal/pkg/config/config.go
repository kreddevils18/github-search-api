package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

type Server struct {
	Port string `json:"port"`
}

type Github struct {
	ApiKey string `json:"apikey"`
}

type DB struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
}
type ElasticSearch struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Address  string `json:"address"`
}

type Organization struct {
	Name string `json:"name"`
	Team string `json:"team"`
}

type Config struct {
	Server        Server        `json:"server"`
	Github        Github        `json:"github"`
	DB            DB            `json:"db"`
	ElasticSearch ElasticSearch `json:"elasticSearch"`
	Organization  Organization  `json:"organization"`
}

func LoadConfig() (*Config, error) {
	var c Config
	var filename string

	if os.Getenv("NODE_ENV") != "production" {
		filename = "development.yaml"
	} else {
		filename = "production.yaml"

	}
	viper.SetConfigFile(fmt.Sprintf("%s/%s", GetPath(), filename))

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&c); err != nil {
		return nil, err
	}

	// load
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("err loading: %v", err)
	}

	// get
	c.Github.ApiKey = os.Getenv("GITHUB_ACCESS_TOKEN")
	c.Organization.Name = os.Getenv("ORGANIZATION")
	c.Organization.Team = os.Getenv("TEAM")

	return &c, nil
}

func GetPath() string {
	path := GetSourcePath() + "/../../../configs///"

	return filepath.Dir(path)
}

func GetSourcePath() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}
