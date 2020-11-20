package learning

import (
	"fmt"
	"log"
	"math/rand"
	"time"
	"xueXi/utils"

	"github.com/trazyn/uiautomator-go"
)

// 将发现的文章保存到结构体中，供点击使用
type cards struct {
	list []generalCard
}
type generalCard struct {
	title    string                // 标题
	position *uiautomator.Position // 坐标
}

// 学习时长，阅读和视频有效阅读时长均为 1分钟
// 此处设置一个阅读最大时间和最小时间，使用其中的随机数值作为阅读时间
var learningTimeMin = 3
var learningTimeMax = 5

// 学习数量
var learningCount = 8

func RandomLearningTime() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(learningTimeMax-learningTimeMin) + learningTimeMin
}

// Learning 学习的方法,视频对应等待(watching())，新闻阅读对应卷动屏幕(reading())
// method选项为news or video
// getCards 文章阅读点击的参照为resourceId，视频也是如此，没有resourceId的文章不去学习
// 根据resourceId获取到所有元素列表，依据元素坐标进行点击操作
func Learning(ua *uiautomator.UIAutomator, method string) error {
	log.Println("start reading:")

	var cards = &uiautomator.Element{}

	// 学习数量
	ln := 0
	for {
		cards = ua.GetElementBySelector(uiautomator.Selector{
			"resourceId": "cn.xuexi.android:id/general_card_title_id", // 新闻的resourceId
		})
		count, err := cards.Count()
		if err != nil {
			return err
		}
		fmt.Println("current page news count:", count)
		for i := 0; i < count; i++ {
			if title, err := cards.Eq(i).GetText(); err == nil {
				fmt.Println("\n", ln+1, title)
			}
			// 对每个card进行点击操作
			err := cards.Eq(i).Click(nil)
			if err != nil {
				return err
			}
			time.Sleep(time.Second)

			switch method {
			case "news":
				err := reading(ua)
				if err != nil {
					return err
				}
			case "video":
				err := watching(ua)
				if err != nil {
					return err
				}
			}

			ln++
			ua.Press("back")
			time.Sleep(time.Second)

			if i+1 == count {
				fmt.Println("swipe screen")
				headPosition, err := utils.GetSelectorPostion(ua, &uiautomator.Selector{
					"resourceId": "cn.xuexi.android:id/comm_head_title",
				})
				if err != nil {
					return err
				}
				if lastP, err := cards.Eq(i).Center(nil); err == nil {
					ua.Swipe(lastP, headPosition, 100)
				}
				time.Sleep(7)
			}

		}
		if ln >= learningCount {
			break
		}

	}

	log.Println("end reading.")

	return nil
}

func reading(ua *uiautomator.UIAutomator) error {
	// 卷动屏幕前稳定1秒
	time.Sleep(time.Second)

	learningTime := RandomLearningTime()
	utils.LearningSwipe(ua, learningTime)
	time.Sleep(time.Second)

	return nil
}

func watching(ua *uiautomator.UIAutomator) error {
	// 创建一个等待阻塞的通道来等待学习时间
	done := make(chan bool)
	learningTime := RandomLearningTime()
	go func() {
		for i := 0; i < learningTime; i++ {
			fmt.Print(".")
			time.Sleep(time.Second)
		}
		done <- true
	}()

	// 等待学习完毕
	<-done
	time.Sleep(time.Second)

	return nil
}
