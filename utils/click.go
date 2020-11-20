package utils

import (
	"fmt"
	"time"
	"xueXi/resource"

	errors "github.com/pkg/errors"
	"github.com/trazyn/uiautomator-go"
)

var err = errors.New("utils/press.go")

// BackHome 返回主页操作
func BackHome(ua *uiautomator.UIAutomator) error {
	// stable
	defer time.Sleep(time.Second)

	appInfo, err := ua.GetCurrentApp()
	if err != nil {
		return err
	}
	if appInfo.Package != resource.AppPackageName {
		err := ua.AppStart(resource.AppPackageName)
		if err != nil {
			return err
		}
	}

	fmt.Println("back to home page")
	position, err := GetSelectorPostion(ua, &uiautomator.Selector{
		"className":  "android.widget.ImageView",
		"resourceId": "cn.xuexi.android:id/home_bottom_tab_icon_large",
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

/*
LearningSwipe 阅读时的滑动方法,跑一个协程监视底部，如果发现底部(此处参照resourceId)，那么滑动结束，返回上一个页面
定义一个学习时长确定满足条件
*/
func LearningSwipe(ua *uiautomator.UIAutomator, learningTime int) {
	// 滑动的起始点和结束点及滑动距离
	pEnd := &uiautomator.Position{
		X: 540,
		Y: 960,
	}
	pStart := &uiautomator.Position{
		X: pEnd.X,
		Y: pEnd.Y + 100,
	}

	// footer的resourceId
	var footer = `xxqg-article-footer`
	var footerTop = make(chan bool)
	var errMsg = make(chan error)
	go func() {
		for {
			se := uiautomator.Selector{
				"resourceId": footer,
			}
			element := ua.GetElementBySelector(se)
			rect, e := element.GetRect()
			// 往下卷动屏幕直到footer距离屏幕顶部小于100
			if e != nil {
				errMsg <- e
				continue
			}
			// 注意此处的610,当文章没有评论时，无法卷动footer至最顶端，故此处设置一个footer顶端的最大距离
			if rect.Top < 610 {
				footerTop <- true
			}
			time.Sleep(time.Second)
		}
	}()

	var learningDuration int

	for {
		ua.Swipe(pStart, pEnd, 10)

		if learningDuration == learningTime {
			return
		}
		learningDuration++

		select {
		case <-errMsg:
			fmt.Print("i")
		case <-footerTop:
			fmt.Println("swipe end")
			return
		default:
			fmt.Print(".")
		}
		time.Sleep(time.Second)
	}

}

// GetSelectorPostion 根据Selector来获取单个元素的坐标
func GetSelectorPostion(ua *uiautomator.UIAutomator, selector *uiautomator.Selector) (p *uiautomator.Position, err error) {

	element := ua.GetElementBySelector(*selector)
	count, err := element.Count()
	if err != nil {
		return nil, err
	}

	if count == 0 {
		return nil, errors.New("not found element by : " + fmt.Sprintf("%v", selector))
	}

	position, err := element.Center(nil)
	if err != nil {
		return nil, err
	}

	return position, nil
}

func Swpie(ua *uiautomator.UIAutomator) error {
	defer time.Sleep(time.Second)
	begin := &uiautomator.Position{
		X: 500,
		Y: 1000,
	}
	end := &uiautomator.Position{
		X: begin.X,
		Y: begin.Y + 500,
	}

	err := ua.Swipe(end, begin, 10)
	if err != nil {
		return err
	}

	return nil
}
