package config

import (
	"os"

	"github.com/spf13/cast"
)

// Config ...
type Config struct {
	SMTPEmail     string
	SMTPEmailPass string
	SMTPHost      string
	SMTPPort      string

	ContextTimeout     int
	Environment        string // develop, staging, production
	SignInKey          string
	OtpTimeout         int // seconds
	AccessTokenTimout  int // MINUTES
	UserServiceHost    string
	UserServicePort    int
	AuthConfigPath     string
	CSVFilePath        string
	PostServiceHost    string
	PostServicePort    int
	CommenterviceHost  string
	CommentServicePort int
	RedisAdres         string
	PostgresHost       string
	PostgresPort       int
	PostgresDatabase   string
	PostgresUser       string
	PostgresPassword   string
	//context timeout in seconds
	CtxTimeout     int
	ConnStr        string
	LogLevel       string
	HTTPPort       string
	JWT_SECRET_KEY string
}

// Load loads environment vars and inflates Config
func Load() Config {
	c := Config{}
	c.AuthConfigPath = cast.ToString(getOrReturnDefault("AUTH_CONFIG_PATH", "./config/auth.conf"))
	c.CSVFilePath = cast.ToString(getOrReturnDefault("CSV_FILE_PATH", "./config/auth.csv"))
	// Email sending
	c.SMTPEmail = cast.ToString(getOrReturnDefault("SMTP_EMAIL", "msrbobo@gmail.com"))
	c.SMTPEmailPass = cast.ToString(getOrReturnDefault("SMTP_EMAIL_PASS", "jqpdwcxjrrrvwqip"))
	c.SMTPHost = cast.ToString(getOrReturnDefault("SMTP_HOST", "smtp.gmail.com"))
	c.SMTPPort = cast.ToString(getOrReturnDefault("SMTP_PORT", "587"))

	c.ContextTimeout = cast.ToInt(getOrReturnDefault("CONTEXT_TIMOUT", 7))

	c.AccessTokenTimout = cast.ToInt(getOrReturnDefault("ACCESS_TOKEN_TIMEOUT", 300))

	c.OtpTimeout = cast.ToInt(getOrReturnDefault("OTP_TIMEOUT", 300))
	c.SignInKey = cast.ToString(getOrReturnDefault("SIGN_IN_KEY", "ASJDKLFJASasdFASE2SD2dafa"))
	c.Environment = cast.ToString(getOrReturnDefault("ENVIRONMENT", "develop"))
	c.JWT_SECRET_KEY = "qjslmdrioekmosiklxeklmrfpo4kmoeqoimk"
	c.LogLevel = cast.ToString(getOrReturnDefault("LOG_LEVEL", "debug"))
	c.HTTPPort = cast.ToString(getOrReturnDefault("HTTP_PORT", ":5053"))
	//userservice
	c.UserServiceHost = cast.ToString(getOrReturnDefault("DOCKER_TEST/user-service_HOST", "user-service"))
	c.UserServicePort = cast.ToInt(getOrReturnDefault("DOCKER_TEST/user-service_PORT", 7778))

	//postservicr
	c.PostServiceHost = cast.ToString(getOrReturnDefault("Post_SERVICE_HOST", "localhost"))
	c.PostServicePort = cast.ToInt(getOrReturnDefault("Post_SERVICE_PORT", 8888))
	c.ConnStr = cast.ToString(getOrReturnDefault("Post_SERVICE", "localhost:8888"))
	//comment
	c.CommenterviceHost = cast.ToString(getOrReturnDefault("Comment_SERVICE_HOST", "localhost"))
	c.CommentServicePort = cast.ToInt(getOrReturnDefault("Comment_SERVICE_PORT", 3333))
	c.RedisAdres = cast.ToString(getOrReturnDefault("REDIS_ADRES", "redis:6379"))
	c.CtxTimeout = cast.ToInt(getOrReturnDefault("CTX_TIMEOUT", 7))

	c.PostgresHost = cast.ToString(getOrReturnDefault("POSTGRES_HOST", "db"))
	c.PostgresPort = cast.ToInt(getOrReturnDefault("POSTGRES_PORT", 5433))
	c.PostgresDatabase = cast.ToString(getOrReturnDefault("POSTGRES_DATABASE", "userdb"))
	c.PostgresUser = cast.ToString(getOrReturnDefault("POSTGRES_USER", "postgres"))
	c.PostgresPassword = cast.ToString(getOrReturnDefault("POSTGRES_PASSWORD", "1234"))

	return c
}

func getOrReturnDefault(key string, defaultValue interface{}) interface{} {
	_, exists := os.LookupEnv(key)
	if exists {
		return os.Getenv(key)
	}

	return defaultValue
}
