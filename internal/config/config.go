package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type KafkaConfig struct {
	Brokers     []string `yaml:"brokers"`
	InputTopic  string   `yaml:"input_topic"`
	OutputTopic string   `yaml:"output_topic"`
	GroupID     string   `yaml:"group_id"`
}

type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type S3Config struct {
	Endpoint   string `yaml:"endpoint"`
	AccessKey  string `yaml:"access_key"`
	SecretKey  string `yaml:"secret_key"`
	UseSSL     bool   `yaml:"use_ssl"`
	BucketName string `yaml:"bucket_name"`
}

type MongoConfig struct {
	URI      string `yaml:"uri"`
	Database string `yaml:"database"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type Config struct {
	Kafka    *KafkaConfig    `yaml:"kafka"`
	Postgres *PostgresConfig `yaml:"postgres"`
	S3       *S3Config       `yaml:"s3"`
	Mongo    *MongoConfig    `yaml:"mongo"`
	Redis    *RedisConfig    `yaml:"redis"`
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
	return &config, nil
}
