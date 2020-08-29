package lib

import (
	"github.com/lithammer/shortuuid"
)

func UUID() string {
	return shortuuid.New()
}
