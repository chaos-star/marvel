package Env

import "strings"

const (
	DeployEnvDebug  = "debug"
	DeployEnvTest   = "test"
	DeployEnvProd   = "prod"
	DeployEnvSimula = "simulation"
)

type Env struct {
	appEnv string
}

func Initialize(env string) *Env {
	env = strings.ToLower(env)
	envInst := &Env{}
	envInst.appEnv = env
	if envInst.appEnv == "" {
		envInst.appEnv = DeployEnvDebug
	}
	switch envInst.appEnv {
	case DeployEnvProd:
		fallthrough
	case DeployEnvTest:
		fallthrough
	case DeployEnvDebug:
		return envInst
	}
	return envInst
}

func (e *Env) GetEnv() string {
	return e.appEnv
}
