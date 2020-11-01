package main

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

type Config struct {
	PublicListenAddr   string
	PublicURL          string
	FrontendURL        string
	AppURL             string
	PostgresURL        string
	JwtSigningKey      string
	StaticFilesPath    string
	SMTPHost           string
	MockSendmail       bool
	FastSpringUsername string
	FastSpringPassword string
	FastSpringHash     string
	FastSpringValidate bool
	PrintConfig        bool
	Development        bool
}

var _configInstance *Config
var _configOnce sync.Once

func GetConfig() *Config {
	_configOnce.Do(func() {
		_configInstance = &Config{}
		_configInstance.ReadConfig()
	})
	return _configInstance
}

func (c *Config) ReadConfig() {
	log.Println("Reading config...")
	c.PublicListenAddr = c._GetEnv("PUBLIC_LISTEN_ADDR", "0.0.0.0:8080")
	c.PublicURL = c._GetEnv("PUBLIC_URL", "http://localhost:8080")
	if c.PublicURL[len(c.PublicURL)-1] != '/' {
		c.PublicURL += "/"
	}
	c.FrontendURL = c._GetEnv("FRONTEND_URL", "http://localhost:3000")
	if c.FrontendURL[len(c.FrontendURL)-1] != '/' {
		c.FrontendURL += "/"
	}
	c.AppURL = c._GetEnv("APP_URL", "exp://localhost:19000")
	if c.AppURL[len(c.AppURL)-1] != '/' {
		c.AppURL += "/"
	}
	c.StaticFilesPath = c._GetEnv("STATIC_FILES_PATH", "/app/adminui")
	if c.StaticFilesPath[len(c.StaticFilesPath)-1] != '/' {
		c.StaticFilesPath += "/"
	}
	c.PostgresURL = c._GetEnv("POSTGRES_URL", "postgres://postgres:root@localhost/flexspace?sslmode=disable")
	c.JwtSigningKey = c._GetEnv("JWT_SIGNING_KEY", "cX32hEwZDCLZ6bCR")
	c.SMTPHost = c._GetEnv("SMTP_HOST", "192.168.40.31:25")
	c.MockSendmail = (c._GetEnv("MOCK_SENDMAIL", "0") == "1")
	c.FastSpringUsername = c._GetEnv("FASTSPRING_USER", "")
	c.FastSpringPassword = c._GetEnv("FASTSPRING_PASS", "")
	c.FastSpringHash = c._GetEnv("FASTSPRING_HASH", "")
	c.FastSpringValidate = (c._GetEnv("FASTSPRING_VALIDATE", "1") == "1")
	c.PrintConfig = (c._GetEnv("PRINT_CONFIG", "0") == "1")
	c.Development = (c._GetEnv("DEV", "0") == "1")
}

func (c *Config) Print() {
	s, _ := json.MarshalIndent(c, "", "\t")
	log.Println("Using config:\n" + string(s))
}

func (c *Config) _GetEnv(key, defaultValue string) string {
	res := os.Getenv(key)
	if res == "" {
		return defaultValue
	}
	return res
}
