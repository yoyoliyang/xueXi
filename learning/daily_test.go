package learning

import (
	"fmt"
	"testing"

	ug "github.com/trazyn/uiautomator-go"
)

func TestEnterDailyAnswers(t *testing.T) {
	ua := ug.New(&ug.Config{
		Host: "192.168.1.52",
		Port: 7912,
	})

	/*
		err := enterDailyAnswers(ua)
		if err != nil {
			t.Fatal(err)
		}
	*/

	questionType, err := getQuestionType(ua)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(questionType)

	/*
		selector := ug.Selector{
			"className": "android.view.View",
		}

		el := ua.GetElementBySelector(selector)
		fmt.Println(el.Count())

	*/
}
