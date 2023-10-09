package Env

import "errors"

const (
	DeployEnvDev  = "Development"
	DeployEnvTest = "Test"
	DeployEnvPre  = "Preview"
	DeployEnvProd = "Production"
)

type Env struct {
	appEnv string
}

func Initialize(env string) (error, *Env) {
	envInst := &Env{}
	envInst.appEnv = env
	if envInst.appEnv == "" {
		envInst.appEnv = DeployEnvDev
	}
	switch envInst.appEnv {
	case DeployEnvProd:
		fallthrough
	case DeployEnvPre:
		fallthrough
	case DeployEnvTest:
		fallthrough
	case DeployEnvDev:
		return nil, envInst
	}
	return errors.New("unknown environment the default will be the development"), nil
}

func (e *Env) GetEnv() string {
	return e.appEnv
}
