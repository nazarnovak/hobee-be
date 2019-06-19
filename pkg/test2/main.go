package test2

import (
	"github.com/nazarnovak/hobee-be/pkg/herrors"
)

type A struct {
	One string
	Two int64
}

func Test2() error {
	return herrors.New("Testing", "1", 2, "3", A{"one", 2})
}
