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
			err := multipleChoice(ua, false)
			if err != nil {
				return err
			}
		case "多选题":
			err := multipleChoice(ua, true)
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

	err = utils.BackHome(ua)
	if err != nil {
		return err
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

func multipleChoice(ua *uiautomator.UIAutomator, method bool) error {
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
			// 获取abcd选项的坐标
			rect, err := element.Eq(i + 1).GetRect()
			if err != nil {
				return err
			}
			// 获取abcd选项文本内容
			_t, err := element.Eq(i + 1).GetText()
			if err != nil {
				return err
			}
			fmt.Println(_t)
			// 判断文本内容是否匹配红色答案
			// 单选时(method=false)，红色答案匹配选项内容，因为部分红色答案相对选项来说是精简的词语
			var re = &regexp.Regexp{}
			if !method {
				re = regexp.MustCompile(fmt.Sprintf("%v{1}", redFontStr))
				if re.MatchString(_t) {
					positions = append(positions, uiautomator.Position{
						X: float32((rect.Left + rect.Right) / 2),
						Y: float32((rect.Bottom + rect.Top) / 2),
					})
				}
			} else {
				re = regexp.MustCompile(_t)
				if re.MatchString(redFontStr) {
					positions = append(positions, uiautomator.Position{
						X: float32((rect.Left + rect.Right) / 2),
						Y: float32((rect.Bottom + rect.Top) / 2),
					})
				}
			}
			fmt.Println(positions)
		}
	}

	for _, p := range positions {
		err := ua.Click(&p)
		if err != nil {
			return err
		}
		time.Sleep(time.Millisecond * 500)
	}

	err = confirm(ua)
	if err != nil {
		return err
	}

	return nil

}

/*
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
*/

// 进入每日答题模块
func enterDailyAnswers(ua *uiautomator.UIAutomator) error {
	defer time.Sleep(time.Second * 2)

	// 进入我的积分页面
	mineElement := ua.GetElementBySelector(uiautomator.Selector{
		"resourceId": "cn.xuexi.android:id/comm_head_xuexi_mine",
	})
	err := mineElement.WaitForExists(1, 5)
	if err != nil {
		return err
	}
	err = mineElement.Click(nil)
	if err != nil {
		return err
	}

	content := [...]string{"我要答题", "每日答题"}

	// 进入我要答题
	fmt.Println("进入", content[0])
	position, err := utils.GetSelectorPostion(ua, &uiautomator.Selector{
		"className":  "android.widget.TextView",
		"resourceId": "cn.xuexi.android:id/user_item_name",
		"package":    "cn.xuexi.android",
		"text":       content[0],
	})
	if err != nil {
		return err
	}

	err = ua.Click(position)
	if err != nil {
		return err
	}

	// 进入每日答题
	fmt.Println("进入", content[1])
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
	err := element.WaitForExists(1, 5)
	if err != nil {
		return err
	}
	err = element.Click(nil)
	if err != nil {
		return errors.New(err.Error() + fmt.Sprint(element))
	}
	return nil
}
