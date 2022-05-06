package Env

import "errors"

const (
	DeployEnvDev  = "Development"
	DeployEnvTest = "Test"
	DeployEnvPre  = "Preview"
	DeployEnvProd = "Production"
)

var (
	appEnv string
)

func Initialize(env string) error{
	appEnv = env
	if env == ""{
		appEnv = DeployEnvDev
	}
	switch env {
	case DeployEnvProd:
		fallthrough
	case DeployEnvPre:
		fallthrough
	case DeployEnvTest:
		fallthrough
	case DeployEnvDev:
		return nil
	}
	return errors.New("unknown environment the default will be the development")
}

func GetEnv() string {
	return appEnv
}
