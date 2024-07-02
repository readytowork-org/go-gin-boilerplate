package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

// Env has environment stored
type Env struct {
	HOST        string `mapstructure:"HOST"`
	TimeZone    string `mapstructure:"TZ"`
	ServerPort  string `mapstructure:"SERVER_PORT"`
	Environment string `mapstructure:"ENVIRONMENT"`
	LogOutput   string `mapstructure:"LOG_OUTPUT"`
	LogLevel    string `mapstructure:"LOG_LEVEL"`

	DBUsername string `mapstructure:"DB_USERNAME"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBName     string `mapstructure:"DB_NAME"`
	SentryDSN  string `mapstructure:"SENTRY_DSN"`

	StorageBucketName string `mapstructure:"STORAGE_BUCKET_NAME"`

	AdminEmail string `mapstructure:"ADMIN_EMAIL"`
	AdminPass  string `mapstructure:"ADMIN_PASS"`
	AdminName  string `mapstructure:"ADMIN_NAME"`

	MailClientID     string `mapstructure:"MAIL_CLIENT_ID"`
	MailClientSecret string `mapstructure:"MAIL_CLIENT_SECRET"`
	MailAccesstoken  string `mapstructure:"MAIL_ACCESS_TOKEN"`
	MailRefreshToken string `mapstructure:"MAIL_REFRESH_TOKEN"`

	AwsS3Region  string `mapstructure:"AWS_S3_REGION"`
	AwsS3Bucket  string `mapstructure:"AWS_S3_BUCKET"`
	AwsAccessKey string `mapstructure:"AWS_ACCESS_KEY"`
	AwsSecretKey string `mapstructure:"AWS_SECRET_KEY"`

	TwilioBaseURL            string `mapstructure:"TWILIO_BASE_URL"`
	TwilioSID                string `mapstructure:"TWILIO_SID"`
	TwilioAuthToken          string `mapstructure:"TWILIO_AUTH_TOKEN"`
	TwilioSMSFrom            string `mapstructure:"TWILIO_SMS_FROM"`
	JwtAccessSecret          string `mapstructure:"JWT_ACCESS_SECRET"`
	JwtRefreshSecret         string `mapstructure:"JWT_REFRESH_SECRET"`
	JwtAccessTokenExpiresAt  int    `mapstructure:"JWT_ACCESS_TOKEN_EXPIRES_AT"`
	JwtRefreshTokenExpiresAt int    `mapstructure:"JWT_REFRESH_TOKEN_EXPIRES_AT"`

	RateLimitPeriod   time.Duration `mapstructure:"RATE_LIMIT_PERIOD"`
	RateLimitRequests int64         `mapstructure:"RATE_LIMIT_REQUESTS"`

	ProjectName       string `mapstructure:"PROJECT_NAME"`
	BillingAccountId  string `mapstructure:"BILLING_ACCOUNT_ID"`
	BudgetDisplayName string `mapstructure:"BUDGET_DISPLAY_NAME"`
	BudgetAmount      int64  `mapstructure:"BUDGET_AMOUNT"`
	SetBudget         int    `mapstructure:"SET_BUDGET"`

	StripeSecretKey   string `mapstructure:"STRIPE_SECRET_KEY"`
	StripeProductID   string `mapstructure:"STRIPE_PRODUCT_ID"`
	StripeWebhookKey  string `mapstructure:"STRIPE_WEBHOOK_KEY"`
	StripeRedirectUrl string `mapstructure:"STRIPE_REDIRECT_URL"`
}

type EnvPath string

func (p EnvPath) ToString() string {
	return string(p)
}

// NewEnv creates a new environment
func NewEnv(envPath EnvPath) Env {
	env := Env{}
	viper.SetConfigFile(envPath.ToString())

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalf("☠️ Env config file not found: %+v", err)
		} else {
			log.Fatalf("☠️ Env config file error: %+v", err)
		}
	}

	if err := viper.Unmarshal(&env); err != nil {
		log.Fatalf("☠️ environment can't be loaded: %+v", err)
	}

	if env.TimeZone == "" {
		env.TimeZone = "UTC"
	}

	return env
}
