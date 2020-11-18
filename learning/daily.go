package learning

import (
	"errors"
	"fmt"
	"time"
	"xueXi/utils"

	"github.com/trazyn/uiautomator-go"
)

//

type question struct {
	inputPosition uiautomator.Position // 填空题输入坐标
	answer        []answer             //选择题列表
}

type answer struct {
	title    string
	position uiautomator.Position
}

func AnswerTheQuestion(ua *uiautomator.UIAutomator) error {
	questionType, err := getQuestionType(ua)
	if err != nil {
		return err
	}
	switch questionType {
	case "单选":

	}

	return nil
}

func singleChocie(ua *uiautomator.UIAutomator) error {

	err := utils.Swpie(ua)
	if err != nil {
		return err
	}

	position, err := getAnswerTheQuestionPostion(ua, &uiautomator.Selector{"className": "android.view.View"}, "查看提示")
	if err != nil {
		return err
	}

	err = ua.Click(position)
	if err != nil {
		return err
	}

	ans, err := utils.GetRedFontString(ua)
	if err != nil {
		return err
	}

	fmt.Println("answer is:", ans)

}

// 进入每日答题模块
func enterDailyAnswers(ua *uiautomator.UIAutomator) error {
	err := utils.BackHome(ua)
	if err != nil {
		return err
	}

	// 进入我的积分页面
	err = utils.ReSourceIDClick(ua, "cn.xuexi.android:id/comm_head_xuexi_mine")
	if err != nil {
		return err
	}

	content := [...]string{"我要答题", "每日答题"}

	selector := &uiautomator.Selector{
		"resourceId": "cn.xuexi.android:id/user_item_name",
	}
	position, err := getAnswerTheQuestionPostion(ua, selector, content[0])
	if err != nil {
		return err
	}
	err = ua.Click(position)
	if err != nil {
		return err
	}
	time.Sleep(time.Second)

	selector = &uiautomator.Selector{
		"className": "android.view.View",
	}
	position, err = getAnswerTheQuestionPostion(ua, selector, content[1])
	if err != nil {
		return err
	}

	err = ua.Click(position)
	if err != nil {
		return err
	}
	defer time.Sleep(time.Second)

	return nil

}

func getAnswerTheQuestionPostion(ua *uiautomator.UIAutomator, selector *uiautomator.Selector, text string) (p *uiautomator.Position, err error) {

	element := ua.GetElementBySelector(*selector)
	count, err := element.Count()
	if err != nil {
		return nil, err
	}

	if count == 0 {
		return nil, errors.New("not found element by : " + fmt.Sprintf("%v", selector))
	}

	for i := 0; i < count; i++ {
		_element := element.Eq(i)
		t, err := _element.GetText()
		if err != nil {
			return nil, err
		}
		fmt.Println(t)
		if t == text {
			rect, err := _element.GetRect()
			if err != nil {
				return nil, err
			}

			p = &uiautomator.Position{
				X: float32((rect.Left + rect.Right) / 2),
				Y: float32((rect.Bottom + rect.Top) / 2),
			}

			return p, nil
		}
	}
	return nil, errors.New("not found element " + text)
}

func getQuestionType(ua *uiautomator.UIAutomator) (string, error) {

	selector := &uiautomator.Selector{
		"className": "android.view.View",
	}

	element := ua.GetElementBySelector(*selector)
	count, err := element.Count()
	if err != nil {
		return "", err
	}

	if count == 0 {
		return "", errors.New("not found element by : " + fmt.Sprintf("%v", selector))
	}

	for i := 0; i < count; i++ {
		_element := element.Eq(i)
		t, err := _element.GetText()
		if err != nil {
			return "", err
		}
		switch t {
		case "单选题":
			return t, nil
		case "多选题":
			return t, nil
		case "填空题":
			return t, nil
		}

	}
	return "", nil
}
