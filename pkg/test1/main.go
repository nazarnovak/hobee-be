package test1

import (
	"hobee-be/pkg/herrors"
	"hobee-be/pkg/test2"
)

func Test1() error {
	return herrors.Wrap(test2.Test2())
}
