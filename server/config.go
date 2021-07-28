package main

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

type Config struct {
	PublicListenAddr    string
	PublicURL           string
	FrontendURL         string
	AppURL              string
	PostgresURL         string
	JwtSigningKey       string
	StaticAdminUiPath   string
	StaticBookingUiPath string
	SMTPHost            string
	MockSendmail        bool
	PrintConfig         bool
	Development         bool
	InitOrgName         string
	InitOrgDomain       string
	InitOrgUser         string
	InitOrgPass         string
	InitOrgCountry      string
	InitOrgLanguage     string
	SignupDomain        string
	SignupAdmin         string
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
	c.Development = (c._GetEnv("DEV", "0") == "1")
	c.PublicListenAddr = c._GetEnv("PUBLIC_LISTEN_ADDR", "0.0.0.0:8080")
	c.PublicURL = c._GetEnv("PUBLIC_URL", "http://localhost:8080")
	if c.PublicURL[len(c.PublicURL)-1] != '/' {
		c.PublicURL += "/"
	}
	if c.Development {
		c.FrontendURL = c._GetEnv("FRONTEND_URL", "http://localhost:3000")
	} else {
		c.FrontendURL = c._GetEnv("FRONTEND_URL", "http://localhost:8080")
	}
	if c.FrontendURL[len(c.FrontendURL)-1] != '/' {
		c.FrontendURL += "/"
	}
	if c.Development {
		c.AppURL = c._GetEnv("APP_URL", "exp://localhost:19000")
	} else {
		c.AppURL = c._GetEnv("APP_URL", "seatsurfing:///")
	}
	if c.AppURL[len(c.AppURL)-1] != '/' {
		c.AppURL += "/"
	}
	c.StaticAdminUiPath = c._GetEnv("STATIC_ADMIN_UI_PATH", "/app/adminui")
	if c.StaticAdminUiPath[len(c.StaticAdminUiPath)-1] != '/' {
		c.StaticAdminUiPath += "/"
	}
	c.StaticBookingUiPath = c._GetEnv("STATIC_BOOKING_UI_PATH", "/app/bookingui")
	if c.StaticBookingUiPath[len(c.StaticBookingUiPath)-1] != '/' {
		c.StaticBookingUiPath += "/"
	}
	c.PostgresURL = c._GetEnv("POSTGRES_URL", "postgres://postgres:root@localhost/seatsurfing?sslmode=disable")
	c.JwtSigningKey = c._GetEnv("JWT_SIGNING_KEY", "cX32hEwZDCLZ6bCR")
	c.SMTPHost = c._GetEnv("SMTP_HOST", "127.0.0.1:25")
	c.MockSendmail = (c._GetEnv("MOCK_SENDMAIL", "0") == "1")
	c.PrintConfig = (c._GetEnv("PRINT_CONFIG", "0") == "1")
	c.InitOrgName = c._GetEnv("INIT_ORG_NAME", "Sample Company")
	c.InitOrgDomain = c._GetEnv("INIT_ORG_DOMAIN", "seatsurfing.de")
	c.InitOrgUser = c._GetEnv("INIT_ORG_USER", "admin")
	c.InitOrgPass = c._GetEnv("INIT_ORG_PASS", "12345678")
	c.InitOrgCountry = c._GetEnv("INIT_ORG_COUNTRY", "DE")
	c.InitOrgLanguage = c._GetEnv("INIT_ORG_LANGUAGE", "de")
	c.SignupDomain = c._GetEnv("SIGNUP_DOMAIN", ".on.seatsurfing.de")
	c.SignupAdmin = c._GetEnv("SIGNUP_ADMIN", "admin")
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
