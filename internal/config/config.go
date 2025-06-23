package config

import (
	"os"
)

type AWSConfig struct {
	AccessKey string
	SecretKey string
	Region    string
	Bucket    string
	Endpoint  string
}

func LoadAWSConfig() AWSConfig {
	return AWSConfig{
		AccessKey: getEnv("AWS_ACCESS_KEY_ID"),
		SecretKey: getEnv("AWS_SECRET_ACCESS_KEY"),
		Region:    getEnv("AWS_REGION"),
		Bucket:    getEnv("AWS_BUCKET"),
		Endpoint:  getEnv("S3_ENDPOINT"),
	}
}

func getEnv(key string) string {
	val := os.Getenv(key)
	return val
}
