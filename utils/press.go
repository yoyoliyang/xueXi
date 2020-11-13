package utils

import (
	"fmt"
	"log"
	"strconv"
	"time"
	"xueXi/resource"

	errors "github.com/pkg/errors"
	"github.com/trazyn/uiautomator-go"
)

var err = errors.New("utils/press.go")

// WaitDuration 等待操作完毕的时间段
var sleepTime = time.Second * 2

// CheckActivity 用来检测当前的activity界面是否与参数中的一致
func CheckActivity(ua *uiautomator.UIAutomator, activity string) (bool, error) {
	fmt.Printf("check activity(%v)\n", activity)
	app, e := ua.GetCurrentApp()
	if e != nil {
		return false, errors.Wrap(err, e.Error())
	}
	if app.Package != resource.AppPackageName {
		panic("当前非学习强国APP界面")
	}

	if app.Activity != activity {
		return false, nil
	}

	return true, nil
}

// BackHome 返回主页操作
func BackHome(ua *uiautomator.UIAutomator) error {
	fmt.Println("back to home page")

	// 启动一个协程来检测界面
	// 此处的通道缓冲值为返回的次数，否则失败
	var res = make(chan bool, 10)
	var errMsg = make(chan error, 10)
	go func() {
		for {
			result, e := CheckActivity(ua, resource.Activity["home"])
			if e != nil {
				errMsg <- e
			}
			if result {
				res <- result
				errMsg <- nil
				break
			}
			errMsg <- nil
			res <- false
			time.Sleep(time.Second)
		}
	}()

	// 循环点击返回按钮，确保回到主页
	for {
		if e := <-errMsg; e != nil {
			log.Fatal(e)
		}
		if result := <-res; result {
			break
		}

		ua.Press("back")
		time.Sleep(time.Second)
	}

	// stable
	time.Sleep(time.Second)

	return nil
}

// ReSourceIDClick 当屏幕上仅存一个resourceId的时候，使用该方法来点击操作
func ReSourceIDClick(ua *uiautomator.UIAutomator, resourceID string) error {
	fmt.Println("resourceIdClick")
	se := uiautomator.Selector{
		"resourceId": resourceID,
	}
	element := ua.GetElementBySelector(se)
	if count, e := element.Count(); e == nil {
		switch count {
		case 0:
			return errors.Wrap(err, "ReSourceIdClick : not found resourceId "+resourceID)
		case 1:
			fmt.Println("found resourceId:", resourceID)
			err = element.Click(nil)
			time.Sleep(sleepTime)
			if err != nil {
				return errors.Wrap(err, "ReSourceIdClick :"+err.Error())
			}
		default:
			return errors.Wrap(err, "ReSourceIdClick : resourceId must be unique, found "+strconv.Itoa(count))
		}
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
		time.Sleep(time.Second)

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
	}

}

// ClickPosition 保存一些固定按钮的坐标（比如要闻等等)
func ClickPosition(ua *uiautomator.UIAutomator, name string) error {
	up := &uiautomator.Position{}
	switch name {
	// 要闻按钮位置
	case "news":
		up = &uiautomator.Position{
			X: 222,
			Y: 255,
		}
	case "tv": // 电视台按钮
		up = &uiautomator.Position{
			X: 672,
			Y: 1746,
		}
		ua.Click(up)
	case "tvNews": // 联播频道
		up = &uiautomator.Position{
			X: 606,
			Y: 255,
		}
	}
	err := ua.Click(up)
	if err != nil {
		return err
	}
	// for stable
	time.Sleep(time.Second)
	return nil
}
