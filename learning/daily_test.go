package learning

import (
	"testing"

	ug "github.com/trazyn/uiautomator-go"
)

var ua = ug.New(&ug.Config{
	Host: "192.168.1.52",
	Port: 7912,
})

func TestEnterDailyAnswers(t *testing.T) {

	/*
		err := enterDailyAnswers(ua)
		if err != nil {
			t.Fatal(err)
		}
	*/

	// fmt.Println(getQuestionType(ua))
	err := AnswerTheQuestion(ua)
	if err != nil {
		t.Fatal(err)
	}

	/*
		selector := ug.Selector{
			"className": "android.view.View",
		}

		el := ua.GetElementBySelector(selector)
		fmt.Println(el.Count())

	*/
}

/*
func TestFillInTheBlanks(t *testing.T) {
	err := fillInTheBlanks(ua)
	if err != nil {
		t.Fatal(err)
	}
}
*/

/*
func TestSingleChoice(t *testing.T) {
	err := singleChoice(ua)
	fmt.Println(err)
}
*/

/*
func TestMultipleChoice(t *testing.T) {
	err := multipleChoice(ua)
	if err != nil {
		t.Fatal(err)
	}
}
*/
