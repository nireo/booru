package lib

import (
	"github.com/lithammer/shortuuid/v3"
)

func UUID() string {
	return shortuuid.New()
}
