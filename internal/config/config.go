package config

import (
	"os"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Config struct {
	DB        DBConfig `yaml:"db"`
	SecretKey string   `yaml:"secretKey"`
}

type DBConfig struct {
	User            string        `yaml:"user"`
	Password        string        `yaml:"password"`
	Host            string        `yaml:"host"`
	DBName          string        `yaml:"dbName"`
	Driver          string        `yaml:"driver"`
	DisableTLS      bool          `yaml:"disableTLS"`
	Timeout         time.Duration `yaml:"timeout"`
	ReadTimeout     time.Duration `yaml:"readTimeout"`
	WriteTimeout    time.Duration `yaml:"writeTimeout"`
	MaxIdleConns    int           `yaml:"maxIdleConns"`
	MaxOpenConns    int           `yaml:"maxOpenConns"`
	ConnMaxIdleTime time.Duration `yaml:"connMaxIdleTime"`
	ConnMaxLifetime time.Duration `yaml:"connMaxLifetime"`
}

func NewConfig(configFilePath string) (Config, error) {
	f, err := os.Open(configFilePath)
	if err != nil {
		return Config{}, errors.Wrap(err, "opening config from file")
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return Config{}, errors.Wrap(err, "decoding yaml config file")
	}

	return cfg, nil
}
