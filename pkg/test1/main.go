package test1

import (
	"github.com/nazarnovak/hobee-be/pkg/herrors"
	"github.com/nazarnovak/hobee-be/pkg/test2"
)

func Test1() error {
	return herrors.Wrap(test2.Test2())
}
