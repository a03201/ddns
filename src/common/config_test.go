package common

import (
	"log"
	"testing"
)

func TestParseConfig(t *testing.T) {
	setting,err:=ParseConfig("../../config.json")
	log.Println(setting.Cloudflare,err)
}
