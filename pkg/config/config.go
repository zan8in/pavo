package config

import (
	"os"
	"path/filepath"

	fileutil "github.com/zan8in/pins/file"
	"gopkg.in/yaml.v3"
)

var (
	defaultConfig = "pavo.yaml"
	configFile    = ""
)

type Config struct {
	Fofa   Fofa   `yaml:"fofa"`
	Hunter Hunter `yaml:"hunter"`
}

type Fofa struct {
	Email string `yaml:"email"`
	Key   string `yaml:"key"`
}

type Hunter struct {
	ApiKey []string `yaml:"api-key"`
}

var (
	HunterApiKeyList     []string
	HunterApiKeyNullList []string
)

func NewConfig() (*Config, error) {
	var err error
	if err = getConfig(); err != nil {
		return nil, err
	}
	return readConfig()
}

func getConfig() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(homeDir, ".config", "pavo")
	if !fileutil.FolderExists(configDir) {
		if err = os.MkdirAll(configDir, 0755); err != nil {
			return err
		}
	}

	configFile = filepath.Join(configDir, defaultConfig)
	if !fileutil.FileExists(configFile) {
		if err = createConfig(configFile); err != nil {
			return err
		}
	}

	return nil
}

func createConfig(configFile string) error {
	configYaml, err := yaml.Marshal(&Config{})
	if err != nil {
		return err
	}

	file, err := os.OpenFile(configFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.Write(configYaml); err != nil {
		return err
	}

	return nil
}

func readConfig() (*Config, error) {
	f, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	config := &Config{}
	if err := yaml.NewDecoder(f).Decode(config); err != nil {
		return nil, err
	}
	return config, nil
}

func (c *Config) IsFofa() bool {
	if len(c.Fofa.Email) == 0 || len(c.Fofa.Key) == 0 {
		return false
	}
	return true
}

func (c *Config) IsHunter() bool {
	HunterApiKeyList = c.Hunter.ApiKey
	return len(c.Hunter.ApiKey) != 0
}

func GetApiKey() string {
	for _, v := range HunterApiKeyList {
		if !IsApiKeyNull(v) {
			return v
		}
	}
	return ""
}

func SetApiKeyNull(apikey string) {
	if !IsApiKeyNull(apikey) {
		HunterApiKeyNullList = append(HunterApiKeyNullList, apikey)
	}
}

func IsApiKeyNull(apikey string) bool {
	for _, v := range HunterApiKeyNullList {
		if v == apikey {
			return true
		}
	}
	return false
}
