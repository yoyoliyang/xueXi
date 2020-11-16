package notice

import (
	"fmt"
	"testing"
)

func TestIftttNotice(t *testing.T) {
	err := IftttNotice("ifttt notice test")
	if err != nil {
		t.Fatal("error for ifttt notice", err)
	}
	fmt.Println("you should get a message from ifttt app")
}
