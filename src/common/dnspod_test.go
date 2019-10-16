package common

import (
	"log"
	"testing"
)

func TestGetCurrentExternalIP(t *testing.T) {
	err:=DnspodRequestByToken("119549","b4e0802e0a92787f383bfbe8e8c69f76","grouk.com","test",
		"1.1.1.1")
	log.Println(err)
}

func TestTes(t *testing.T) {
	err:=TestDnspodRequestByToken("119549","b4e0802e0a92787f383bfbe8e8c69f76")
	t.Log(err)
}