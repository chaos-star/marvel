package Config

import (
	"os"
)

type Loading interface {
	Load() (error, []byte)
}

type DefaultLoading struct {
}

func (l *DefaultLoading) Load() (err error, body []byte) {
	var path = "./Conf/app.toml"
	body, err = os.ReadFile(path)
	if err != nil {
		return
	}
	return
}
