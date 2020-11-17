package utils

import (
	"fmt"
	"testing"
)

func TestGetRedAnswer(t *testing.T) {

	rgb := &Red{}
	rgb, err := GetRedFontImage()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(rgb)

}
