package utils

import (
	"testing"

	ug "github.com/trazyn/uiautomator-go"
)

func TestBackHome(t *testing.T) {
	ua := ug.New(&ug.Config{
		Host: "192.168.1.52",
		Port: 7912,
	})
	err := BackHome(ua)
	if err != nil {
		t.Fatal(err)
	}
}
