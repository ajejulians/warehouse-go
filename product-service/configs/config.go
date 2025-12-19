package configs

import "github.com/spf13/viper"

type App struct {
	AppPort string `json:"app_port"`
	AppEnv  string `json:"app_env"`

	UrlMerchantService string `json:"url_merchant_service"`
	UrlProductService  string `json:"url_product_service"`
}

type SqlDB struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User	 string `json:"user"`
	Password string `json:"password"`
	DBname   string `json:"db_name"`
	DBMaxOpenConns	int    `json:"db_max_open_conns"`
	DBMaxIdleConns	int    `json:"db_max_idle_conns"`
}

type Redis struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
}

type RabbitMQ struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Supabase struct {
	Url string `json:"url"`
	Key string `json:"key"`
	Bucket string `json:"bucket"`
}

type Config struct {
	App      App      `json:"app"`
	SqlDB    SqlDB    `json:"sql_db"`
	Redis    Redis    `json:"redis"`
	RabitMQ  RabbitMQ  `json:"rabitmq"`
	Supabase Supabase `json:"supabase"` 
}

func NewConfig() *Config {
	return &Config{
		App: App{
			AppPort: viper.GetString("APP_PORT"),
			AppEnv:  viper.GetString("APP_ENV"),
			UrlMerchantService: viper.GetString("URL_MERCHANT_SERVICE"),
			UrlProductService: viper.GetString("URL_PRODUCT_SERVICE"),
	},
		SqlDB: SqlDB {
			Host:     viper.GetString("DATABASE_HOST"),
			Port:     viper.GetString("DATABASE_PORT"),
			User:     viper.GetString("DATABASE_USER"),
			Password: viper.GetString("DATABASE_PASSWORD"),
			DBname:   viper.GetString("DATABASE_NAME"),
			DBMaxIdleConns: viper.GetInt("DATABASE_MAX_IDLE_CONNECTION"),
			DBMaxOpenConns: viper.GetInt("DATABASE_MAX_OPEN_CONNECTION"),

	},
		RabitMQ: RabbitMQ{
			Host:    	viper.GetString("RABBITMQ_HOST"),
			Port:    	viper.GetString("RABBITMQ_PORT"),
			Username: 	viper.GetString("RABBITMQ_USERNAME"),
			Password: 	viper.GetString("RABBITMQ_PASSWORD"),
	},
		Redis: Redis{
			Host: viper.GetString("REDIS_HOST"),
			Port: viper.GetString("REDIS_PORT"),
	},
		Supabase: Supabase{
			Url: viper.GetString("SUPABASE_URL"),
			Key: viper.GetString("SUPABASE_KEY"),
			Bucket: viper.GetString("SUPABASE_BUCKET"),
	},
  }
}