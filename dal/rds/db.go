package rds

import "github.com/wymli/relaxadmin/common/config"

type Config struct {
	Username          string
	Password          string
	Host              string
	Port              string
	DB                string
	LogFilePath       string
	WithConsoleLogger bool
	WithFileLogger    bool
}

func Init(c *config.DBConfig) {

}
