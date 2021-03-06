package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"sync"
)

type Config struct {
	PublicListenAddr       string
	PublicURL              string
	FrontendURL            string
	AppURL                 string
	PostgresURL            string
	JwtSigningKey          string
	StaticAdminUiPath      string
	StaticBookingUiPath    string
	SMTPHost               string
	SMTPPort               int
	SMTPSenderAddress      string
	SMTPStartTLS           bool
	SMTPInsecureSkipVerify bool
	SMTPAuth               bool
	SMTPAuthUser           string
	SMTPAuthPass           string
	MockSendmail           bool
	PrintConfig            bool
	Development            bool
	InitOrgName            string
	InitOrgDomain          string
	InitOrgUser            string
	InitOrgPass            string
	InitOrgCountry         string
	InitOrgLanguage        string
	OrgSignupEnabled       bool
	OrgSignupDomain        string
	OrgSignupAdmin         string
	OrgSignupMaxUsers      int
	OrgSignupDelete        bool
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
	c.SMTPHost = c._GetEnv("SMTP_HOST", "127.0.0.1")
	smtpPort, err := strconv.Atoi(c._GetEnv("SMTP_PORT", "25"))
	if err != nil {
		log.Fatal("Could not parse SMTP_PORT to int")
	}
	c.SMTPPort = smtpPort
	c.SMTPStartTLS = (c._GetEnv("SMTP_START_TLS", "0") == "1")
	c.SMTPInsecureSkipVerify = (c._GetEnv("SMTP_INSECURE_SKIP_VERIFY", "0") == "1")
	c.SMTPAuth = (c._GetEnv("SMTP_AUTH", "0") == "1")
	c.SMTPAuthUser = c._GetEnv("SMTP_AUTH_USER", "")
	c.SMTPAuthPass = c._GetEnv("SMTP_AUTH_PASS", "")
	c.SMTPSenderAddress = c._GetEnv("SMTP_SENDER_ADDRESS", "no-reply@seatsurfing.local")
	c.MockSendmail = (c._GetEnv("MOCK_SENDMAIL", "0") == "1")
	c.PrintConfig = (c._GetEnv("PRINT_CONFIG", "0") == "1")
	c.InitOrgName = c._GetEnv("INIT_ORG_NAME", "Sample Company")
	c.InitOrgDomain = c._GetEnv("INIT_ORG_DOMAIN", "seatsurfing.local")
	c.InitOrgUser = c._GetEnv("INIT_ORG_USER", "admin")
	c.InitOrgPass = c._GetEnv("INIT_ORG_PASS", "12345678")
	c.InitOrgCountry = c._GetEnv("INIT_ORG_COUNTRY", "DE")
	c.InitOrgLanguage = c._GetEnv("INIT_ORG_LANGUAGE", "de")
	c.OrgSignupEnabled = (c._GetEnv("ORG_SIGNUP_ENABLED", "0") == "1")
	c.OrgSignupDomain = c._GetEnv("ORG_SIGNUP_DOMAIN", ".on.seatsurfing.local")
	c.OrgSignupAdmin = c._GetEnv("ORG_SIGNUP_ADMIN", "admin")
	maxUsers, err := strconv.Atoi(c._GetEnv("ORG_SIGNUP_MAX_USERS", "10"))
	if err != nil {
		log.Fatal("Could not parse ORG_SIGNUP_MAX_USERS to int")
	}
	c.OrgSignupMaxUsers = maxUsers
	c.OrgSignupDelete = (c._GetEnv("ORG_SIGNUP_DELETE", "0") == "1")
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
