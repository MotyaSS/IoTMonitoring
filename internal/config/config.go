package config

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

type GRPCConfig struct {
	Addr string `yaml:"addr"`
}

type KafkaConfig struct {
	Brokers     []string `yaml:"brokers" validate:"required"`
	InputTopic  *string  `yaml:"input_topic"`
	OutputTopic *string  `yaml:"output_topic"`
	Partition   *int     `yaml:"partition"`
	GroupID     string   `yaml:"group_id" validate:"required"`
}

type PostgresConfig struct {
	Host     string `yaml:"host" validate:"required"`
	Port     string `yaml:"port" validate:"required"`
	Username string `yaml:"username" validate:"required"`
	Password string `yaml:"password" validate:"required"`
	DB       string `yaml:"db" validate:"required"`
}

type S3Config struct {
	Endpoint   string `yaml:"endpoint" validate:"required"`
	AccessKey  string `yaml:"access_key" validate:"required"`
	SecretKey  string `yaml:"secret_key" validate:"required"`
	UseSSL     bool   `yaml:"use_ssl"`
	BucketName string `yaml:"bucket_name" validate:"required"`
}

type MongoConfig struct {
	URI string `yaml:"uri" validate:"required"`
	DB  string `yaml:"db" validate:"required"`
}

type Config struct {
	GRPC     *GRPCConfig     `yaml:"grpc"`
	Kafka    *KafkaConfig    `yaml:"kafka"`
	Postgres *PostgresConfig `yaml:"postgres"`
	S3       *S3Config       `yaml:"s3"`
	Mongo    *MongoConfig    `yaml:"mongo"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	validate := validator.New()
	if err := validate.Struct(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}
	return &config, nil
}
