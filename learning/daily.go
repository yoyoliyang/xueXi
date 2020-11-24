package learning

import (
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"
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

	err := enterDailyAnswers(ua)
	if err != nil {
		return err
	}

	for i := 0; i < 5; i++ {
		questionType, err := getQuestionType(ua)
		fmt.Println(questionType)
		if err != nil {
			return err
		}
		switch questionType {
		case "单选题":
			err := singleChoice(ua)
			if err != nil {
				return err
			}
		case "多选题":
			err := multipleChoice(ua)
			if err != nil {
				return err
			}
		case "填空题":
			err := fillInTheBlanks(ua)
			if err != nil {
				return err
			}
		}

	}

	return nil
}

// 填空题
func fillInTheBlanks(ua *uiautomator.UIAutomator) error {
	defer time.Sleep(time.Second * 2)
	redFontStr, err := getRedAnswer(ua)
	if err != nil {
		return err
	}

	element := ua.GetElementBySelector(uiautomator.Selector{"className": "android.view.View"})
	count, err := element.Count()
	if err != nil {
		return err
	}

	for i := 0; i < count; i++ {
		fmt.Print(".")
		if rect, err := element.Eq(i).GetRect(); err == nil {
			if ((rect.Right-rect.Left) == 87 || (rect.Right-rect.Left) == 84) && (rect.Bottom-rect.Top) == 87 {
				fmt.Println(rect)
				err := ua.SetFastinputIME(true)
				if err != nil {
					return err
				}
				err = ua.Click(&uiautomator.Position{
					Y: float32((rect.Bottom + rect.Top) / 2),
					X: float32((rect.Left + rect.Right) / 2),
				})
				if err != nil {
					return err
				}

				str := base64.StdEncoding.EncodeToString([]byte(redFontStr))
				// adb shell广播模式发送中文字符
				_, err = ua.Shell([]string{"am", "broadcast", "-a", "ADB_INPUT_TEXT", "--es", "text", str}, 5)
				if err != nil {
					return err
				}
				time.Sleep(time.Second)
				fmt.Println("done")
				break
			}
		} else {
			return err
		}
	}

	position, err := utils.GetSelectorPostion(ua, &uiautomator.Selector{
		"className": "android.view.View",
		"text":      "确定",
	})
	if err != nil {
		return err
	}

	err = ua.Click(position)
	if err != nil {
		return err
	}

	return nil
}

func getRedAnswer(ua *uiautomator.UIAutomator) (string, error) {
	err := utils.Swpie(ua)
	if err != nil {
		return "", err
	}

	position, err := utils.GetSelectorPostion(ua, &uiautomator.Selector{
		"className": "android.view.View",
		"text":      "查看提示",
	})
	if err != nil {
		return "", err
	}

	err = ua.Click(position)
	if err != nil {
		return "", err
	}

	time.Sleep(time.Second * 2)
	redFontStr, err := utils.GetRedFontString(ua)
	if err != nil {
		return "", err
	}
	ua.Press("back")
	time.Sleep(time.Second)

	fmt.Println("red answer:", redFontStr)
	return redFontStr, nil
}

func multipleChoice(ua *uiautomator.UIAutomator) error {
	defer time.Sleep(time.Second * 2)

	redFontStr, err := getRedAnswer(ua)
	if err != nil {
		return err
	}

	// 创建一个多选题的答案坐标切片
	positions := make([]uiautomator.Position, 0)
	element := ua.GetElementBySelector(uiautomator.Selector{"className": "android.view.View"})
	count, err := element.Count()
	if err != nil {
		return err
	}

	for i := 0; i < count; i++ {
		text, err := element.Eq(i).GetText()
		if err != nil {
			return err
		}

		// 根据ABCD四个字母的位置来确定答案的位置和内容
		if text == "A." || text == "B." || text == "C." || text == "D." {
			fmt.Println(positions)
			if rect, err := element.Eq(i + 1).GetRect(); err == nil {
				if _t, err := element.Eq(i + 1).GetText(); err == nil {
					fmt.Println(_t)
					if m, err := regexp.MatchString(_t, redFontStr); err == nil {
						if m {
							positions = append(positions, uiautomator.Position{
								X: float32((rect.Left + rect.Right) / 2),
								Y: float32((rect.Bottom + rect.Top) / 2),
							})
						}

					} else {
						return nil
					}
				}
			}
		}
	}

	for _, p := range positions {
		err := ua.Click(&p)
		if err != nil {
			return err
		}
		time.Sleep(time.Second)
	}

	err = confirm(ua)
	if err != nil {
		return err
	}

	return nil

}

// 单选题
func singleChoice(ua *uiautomator.UIAutomator) error {
	defer time.Sleep(time.Second * 2)

	redFontStr, err := getRedAnswer(ua)
	if err != nil {
		return err
	}

	position, err := utils.GetSelectorPostion(ua, &uiautomator.Selector{
		"className": "android.view.View",
		"text":      redFontStr,
	})
	if err != nil {
		return err
	}

	err = ua.Click(position)
	if err != nil {
		return err
	}

	err = confirm(ua)
	if err != nil {
		return err
	}

	return nil

}

// 进入每日答题模块
func enterDailyAnswers(ua *uiautomator.UIAutomator) error {
	defer time.Sleep(time.Second * 2)

	// 进入我的积分页面
	mineElement := ua.GetElementBySelector(uiautomator.Selector{
		"resourceId": "cn.xuexi.android:id/comm_head_xuexi_mine",
	})
	err := mineElement.Click(nil)
	if err != nil {
		return err
	}
	scoreValue := ua.GetElementBySelector(uiautomator.Selector{
		"resourceId": "cn.xuexi.android:id/my_score_value",
		"className":  "android.widget.TextView",
	})
	scoreValue.WaitForExists(1, 5)
	currentScore, err := scoreValue.GetText()
	fmt.Println("current scores: ", currentScore)

	content := [...]string{"我要答题", "每日答题"}

	// 进入我要答题
	position, err := utils.GetSelectorPostion(ua, &uiautomator.Selector{
		"resourceId": "cn.xuexi.android:id/user_item_name",
		"text":       content[0],
	})
	if err != nil {
		return err
	}

	err = ua.Click(position)
	if err != nil {
		return err
	}
	time.Sleep(time.Second)

	// 进入每日答题
	position, err = utils.GetSelectorPostion(ua, &uiautomator.Selector{
		"className": "android.view.View",
		"package":   "cn.xuexi.android",
		"text":      content[1],
	})
	if err != nil {
		return err
	}

	err = ua.Click(position)
	if err != nil {
		return err
	}

	return nil

}

func getQuestionType(ua *uiautomator.UIAutomator) (string, error) {

	selector := uiautomator.Selector{
		"className": "android.view.View",
	}

	titles := [...]string{"单选题", "多选题", "填空题"}

	for _, t := range titles {
		selector["text"] = t
		element := ua.GetElementBySelector(selector)
		if qType, err := element.GetText(); err == nil {
			return qType, nil
		}
	}

	return "", errors.New("not found element by : " + fmt.Sprintf("%v", selector))
}

// 点击“确定”按钮
func confirm(ua *uiautomator.UIAutomator) error {
	element := ua.GetElementBySelector(uiautomator.Selector{
		"className": "android.view.View",
		"text":      "确定",
		"clickable": true,
	})
	fmt.Println(element.Count())
	err := element.Click(nil)
	if err != nil {
		return errors.New(err.Error() + fmt.Sprint(element))
	}
	return nil
}
