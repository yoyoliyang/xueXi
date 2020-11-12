package score

import (
	"fmt"
	"xueXi/resource"
	"xueXi/utils"

	"github.com/pkg/errors"
	"github.com/trazyn/uiautomator-go"
)

type Score struct {
	All      int // 所有积分
	Login    int // 登录
	Read     int // 选读文章
	AudioV   int // 视听学习
	AudioL   int // 视听学习时长
	DailyAns int // 每日答题
}

var err = errors.New("score/points.go")

// GetScore 返回用户当前的学习积分结构体
func GetScore(ua *uiautomator.UIAutomator) error {
	fmt.Print("获取当前学习积分 > ")
	// 返回个人积分页面
	e := func() error {
		b, e := utils.CheckActivity(ua, resource.Activity["learnScore"])
		if e != nil {
			return errors.Wrap(err, e.Error())
		}
		if !b {
			e := utils.BackHome(ua)
			if e != nil {
				return errors.Wrap(err, e.Error())
			}
			e = utils.ReSourceIDClick(ua, "cn.xuexi.android:id/comm_head_xuexi_score")
			if e != nil {
				return errors.Wrap(err, e.Error())
			}
		}
		return nil
	}()
	if e != nil {
		return errors.Wrap(err, e.Error())
	}

	scoreSe := uiautomator.Selector{
		"className": "android.view.View",
	}

	element := ua.GetElementBySelector(scoreSe)
	count, e := element.Count()
	if e != nil {
		return errors.Wrap(err, e.Error())
	}
	for i := 0; i < count; i++ {
		el := element.Eq(i)
		fmt.Println(el.GetText())
	}
	return nil
}
