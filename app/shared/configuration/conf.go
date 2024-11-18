package configuration

import (
	"secret-management/app/shared/constants"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type Conf struct {
	envLoader         EnvLoader
	PORT              string `required:"true"`
	VERSION           string `required:"true"`
	ENVIRONMENT       string `required:"false"`
	PROJECT_NAME      string `required:"false"`
	GOOGLE_PROJECT_ID string `required:"false"`
}

func init() {
	ioc.Registry(NewConf, NewEnvLoader)
}
func NewConf(env EnvLoader) (Conf, error) {
	conf := Conf{
		envLoader:         env,
		PORT:              env.Get("PORT"),
		VERSION:           env.Get(constants.Version),
		ENVIRONMENT:       env.Get("ENVIRONMENT"),
		PROJECT_NAME:      env.Get("PROJECT_NAME"),
		GOOGLE_PROJECT_ID: env.Get("GOOGLE_PROJECT_ID"),
	}
	if conf.PORT == "" {
		conf.PORT = "8080"
	}
	return validateConfig(conf)
}

func (c Conf) LoadFromSystem(key string) string {
	return c.envLoader.Get(key)
}
