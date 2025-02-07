package env

import "github.com/Netflix/go-env"

type environment struct {
	APIPort           string  `env:"API_PORT,default=8080"`
	APIName           string  `env:"API_NAME,default=investment-analyzer"`
	APIVersion        string  `env:"API_VERSION"`
	DBName            string  `env:"DB_NAME"`
	DBUser            string  `env:"DB_USER"`
	DBPassword        string  `env:"DB_PASSWORD"`
	DBHost            string  `env:"DB_HOST"`
	DBPort            string  `env:"DB_PORT"`
	SplunkHost        string  `env:"SPLUNK_HOST"`
	SplunkToken       string  `env:"SPLUNK_TOKEN"`
	SplunkAssync      bool    `env:"SPLUNK_ASSYNC,default=false"`
	CacheHost         string  `env:"CACHE_HOST"`
	CachePort         string  `env:"CACHE_PORT"`
	CachePassword     string  `env:"CACHE_PASSWORD"`
	CacheReadTimeout  int64   `env:"CACHE_READ_TIMEOUT"`
	CacheWriteTimeout int64   `env:"CACHE_WRITE_TIMEOUT"`
	AllowOrigins      string  `env:"ALLOW_ORIGINS"`
	Environment       string  `env:"ENVIRONMENT"`
	CashInHand        float64 `env:"CashInHand,default=100000.00"`
}

var Env environment

func LoadEnv() (err error) {
	_, err = env.UnmarshalFromEnviron(&Env)
	return
}
