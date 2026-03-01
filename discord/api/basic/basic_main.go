package basic

import (
	"sync"
)

var instance *apiInterface
var once sync.Once

type apiInterface struct {
}

func GetInstance() *apiInterface {
	once.Do(func() {
		instance = &apiInterface{}
		instance.Init()
	})
	return instance
}

func (a *apiInterface) Init() {
}
