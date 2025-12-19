package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

type App struct {
	AppPort string `json:"app_port"`
	AppEnv  string `json:"app_env"`
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

type Email struct {
	Host		string `json:"host"`
	Port		int		`json:"port"`
	User		string `json:"user"`
	Password	string `json:"password"`
	From		string `json:"from"`
}

type Config struct {
	App      App      	`json:"app"`
	SqlDB    SqlDB    	`json:"sql_db"`
	Redis    Redis    	`json:"redis"`
	RabitMQ  RabbitMQ  	`json:"rabitmq"`
	Email	 Email	 	`json:"email"`
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
			Username: 	viper.GetString("RABBITMQ_USER"),
			Password: 	viper.GetString("RABBITMQ_PASSWORD"),
	},
		Redis: Redis{
			Host: viper.GetString("REDIS_HOST"),
			Port: viper.GetString("REDIS_PORT"),
	},
		Email: Email{
			Host: viper.GetString("EMAIL_HOST"),
			Port: viper.GetInt("EMAIL_PORT"),
			User: viper.GetString("EMAIL_USER"),
			Password: viper.GetString("EMAIL_PASSWROD"),
			From: viper.GetString("EMAIL_FROM"),
	},
  }
}