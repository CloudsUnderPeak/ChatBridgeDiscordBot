package test

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
	})
	return instance
}
