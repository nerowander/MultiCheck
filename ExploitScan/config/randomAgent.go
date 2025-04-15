package config

import (
	"github.com/corpix/uarand"
)

func RandomUserAgent() string {
	randomUA := uarand.GetRandom()
	return randomUA
	//log.Printf("Random: %s", random)
}
