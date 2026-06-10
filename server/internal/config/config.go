package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort   string
	DBHost       string
	DBPort       string
	DBUser       string
	DBPassword   string
	DBName       string
	JWTSecret         string
	JWTPrivateKeyPath string
	JWTPublicKeyPath  string
	DefaultAdmin      string
	DefaultPass  string
	WebOrigin    string
	AdminOrigin  string
	CORSOrigins  string
	GinMode      string
	UploadDir    string
	DeepSeekKey  string
	DeepSeekURL  string
}

func Load() Config {
	loadDotEnv()

	cfg := Config{
		ServerPort:   getEnv("SERVER_PORT", "8080"),
		DBHost:       getEnv("DB_HOST", "127.0.0.1"),
		DBPort:       getEnv("DB_PORT", "3306"),
		DBUser:       getEnv("DB_USER", "root"),
		DBPassword:   getEnv("DB_PASSWORD", ""),
		DBName:       getEnv("DB_NAME", "blog_system"),
		JWTSecret:         getEnv("JWT_SECRET", ""),
		JWTPrivateKeyPath: getEnv("JWT_PRIVATE_KEY_PATH", "./keys/private.pem"),
		JWTPublicKeyPath:  getEnv("JWT_PUBLIC_KEY_PATH", "./keys/public.pem"),
		DefaultAdmin:      getEnv("DEFAULT_ADMIN_USERNAME", "admin"),
		DefaultPass:  getEnv("DEFAULT_ADMIN_PASSWORD", "admin123"),
		WebOrigin:    getEnv("WEB_ORIGIN", "http://localhost:5173"),
		AdminOrigin:  getEnv("ADMIN_ORIGIN", "http://localhost:5174"),
		CORSOrigins:  getEnv("CORS_ORIGINS", ""),
		GinMode:      getEnv("GIN_MODE", "release"),
		UploadDir:    getEnv("UPLOAD_DIR", "./uploads"),
		DeepSeekKey:  getEnv("DEEPSEEK_API_KEY", ""),
		DeepSeekURL:  getEnv("DEEPSEEK_BASE_URL", "https://api.deepseek.com/v1"),
	}
	cfg.normalize()
	return cfg
}

func loadDotEnv() {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "load .env: %v\n", err)
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func (cfg *Config) normalize() {
	cfg.ServerPort = strings.TrimSpace(cfg.ServerPort)
	cfg.DBHost = strings.TrimSpace(cfg.DBHost)
	cfg.DBPort = strings.TrimSpace(cfg.DBPort)
	cfg.DBUser = strings.TrimSpace(cfg.DBUser)
	cfg.DBName = strings.TrimSpace(cfg.DBName)
	cfg.JWTSecret = strings.TrimSpace(cfg.JWTSecret)
	cfg.DefaultAdmin = strings.TrimSpace(cfg.DefaultAdmin)
	cfg.DefaultPass = strings.TrimSpace(cfg.DefaultPass)
	cfg.WebOrigin = strings.TrimSpace(cfg.WebOrigin)
	cfg.AdminOrigin = strings.TrimSpace(cfg.AdminOrigin)
	cfg.CORSOrigins = strings.TrimSpace(cfg.CORSOrigins)
	cfg.GinMode = strings.TrimSpace(cfg.GinMode)
	cfg.UploadDir = strings.TrimSpace(cfg.UploadDir)
	cfg.DeepSeekKey = strings.TrimSpace(cfg.DeepSeekKey)
	cfg.DeepSeekURL = strings.TrimSpace(cfg.DeepSeekURL)
	if cfg.DeepSeekURL == "" {
		cfg.DeepSeekURL = "https://api.deepseek.com/v1"
	}
}

func (cfg Config) Validate() error {
	required := []struct {
		name  string
		value string
	}{
		{name: "SERVER_PORT", value: cfg.ServerPort},
		{name: "DB_HOST", value: cfg.DBHost},
		{name: "DB_PORT", value: cfg.DBPort},
		{name: "DB_USER", value: cfg.DBUser},
		{name: "DB_NAME", value: cfg.DBName},
		{name: "DEFAULT_ADMIN_USERNAME", value: cfg.DefaultAdmin},
		{name: "DEFAULT_ADMIN_PASSWORD", value: cfg.DefaultPass},
		{name: "WEB_ORIGIN", value: cfg.WebOrigin},
		{name: "ADMIN_ORIGIN", value: cfg.AdminOrigin},
		{name: "GIN_MODE", value: cfg.GinMode},
		{name: "UPLOAD_DIR", value: cfg.UploadDir},
	}

	for _, item := range required {
		if strings.TrimSpace(item.value) == "" {
			return fmt.Errorf("配置项 %s 不能为空", item.name)
		}
	}

	return nil
}

func (cfg Config) GetDeepSeekKey() string { return cfg.DeepSeekKey }
func (cfg Config) GetDeepSeekURL() string { return cfg.DeepSeekURL }

func (cfg Config) StartupSummary() string {
	return fmt.Sprintf(
		"port=%s db=%s@%s:%s/%s web=%s admin=%s mode=%s uploads=%s",
		cfg.ServerPort,
		cfg.DBUser,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		cfg.WebOrigin,
		cfg.AdminOrigin,
		cfg.GinMode,
		cfg.UploadDir,
	)
}
