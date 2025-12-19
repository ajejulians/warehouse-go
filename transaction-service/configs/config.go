package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

type App struct {
	AppPort string `json:"app_port"`
	AppEnv  string `json:"app_env"`

	UrlProductService	string 		`json:"url_product_service"`
	UrlUserService 		string	 	`json:"url_user_service"`
	UrlMerchantService string 		`json:"url_merchant_service"`
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

type Midtrans struct {
	ServerKey		string 	`json:"server_key"`
	ClientKey		string 	`json:"client_key"`
	MerchantID 		string	`json:"merchant_id"`
	IsProduction 	bool 	`json:"is_production"`
}

type Config struct {
	App       App      	`json:"app"`
	SqlDB     SqlDB    	`json:"sql_db"`
	Redis     Redis    	`json:"redis"`
	RabbitMQ  RabbitMQ  `json:"rabitmq"`
	Supabase  Supabase 	`json:"supabase"` 
	Midtrans  Midtrans  `json:"midtrans"`
}

//URL Returns the RabbitMQ connection string
func (r *RabbitMQ) URL() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s/", r.Username, r.Password, r.Host, r.Port)
}
 

func NewConfig() *Config {
	return &Config{
		App: App{
			AppPort: viper.GetString("APP_PORT"),
			AppEnv:  viper.GetString("APP_ENV"),
			UrlProductService: viper.GetString("URL_PRODUCT_SERVICE"),
			UrlUserService: viper.GetString("URL_USER_SERVICE"),
			UrlMerchantService: viper.GetString("URL_MERCHANT_SERVICE"),
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
		RabbitMQ: RabbitMQ{
			Host:    	viper.GetString("RABBITMQ_HOST"),
			Port:    	viper.GetString("RABBITMQ_PORT"),
			Username: 	viper.GetString("RABBITMQ_USER"),
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
		Midtrans: Midtrans{
			ServerKey: viper.GetString("MIDTRANS_SERVER_KEY"),
			ClientKey: viper.GetString("MIDTRANS_CLIENT_KEY"),
			MerchantID: viper.GetString("MIDTRANS_MERCHANT_ID"),
			IsProduction: viper.GetBool("MIDTRANS_IS_PRODUCTION"),
		},
  }
}