package config

import "errors"

type Config struct {
	env          string
	mongoURI     string
	databaseName string
	redisURI     string
}

func NewConfig() *Config {
	return &Config{}
}

func (config *Config) SetEnv(env string) {
	config.env = env
}

func (config *Config) SetMongoURI(mongoURI string) {
	config.mongoURI = mongoURI
}

func (config *Config) SetDatabaseName(databaseName string) {
	config.databaseName = databaseName
}

func (config *Config) SetRedisURI(redisURI string) {
	config.redisURI = redisURI
}

func (config *Config) GetEnv() string {
	return config.env
}

func (config *Config) GetMongoURI() string {
	return config.mongoURI
}

func (config *Config) GetDatabaseName() string {
	return config.databaseName
}

func (config *Config) GetRedisURI() string {
	return config.redisURI
}

func (config *Config) Validate() error {
	if config.env == "" {
		return errors.New("GO_ENV is not set")
	}
	if config.mongoURI == "" {
		return errors.New("MONGO_URI is not set")
	}
	if config.databaseName == "" {
		return errors.New("MONGODB_DATABASE is not set")
	}
	if config.redisURI == "" {
		return errors.New("REDIS_URI is not set")
	}
	return nil
}
