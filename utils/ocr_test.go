package utils

import (
	"fmt"
	"testing"

	ug "github.com/trazyn/uiautomator-go"
)

func TestPrintAnswer(t *testing.T) {
	ua := ug.New(&ug.Config{
		Host: "192.168.1.52",
		Port: 7912,
	})

	str, err := GetRedFontString(ua)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(str)

}
