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
var learningTimeMin = 61
var learningTimeMax = 71

// 学习数量
var learningCount = 8

func RandomLearningTime() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(learningTimeMax-learningTimeMin) + learningTimeMin
}

// Learning 学习的方法,视频对应等待(watching())，新闻阅读对应卷动屏幕(reading())
// method选项为news or video
func Learning(ua *uiautomator.UIAutomator, method string) error {
	log.Println("starting reading:")

	cards := &cards{}
	cards, err := cards.GetCards(ua)
	if err != nil {
		return err
	}

	ln := 0
	for {
		for _, card := range cards.list {
			fmt.Println(ln, card.title)
			// 点击新闻标题进入
			err := ua.Click(card.position)
			if err != nil {
				return err
			}

			switch method {
			case "news":
				reading(ua)
			case "video":
				watching()
			}

			// 返回首页
			err = utils.BackHome(ua)
			if err != nil {
				return err
			}

			ln++
			if ln >= learningCount {
				break
			}
		}
		if ln >= learningCount {
			break
		}
		nc, err := cards.cardSwipe(ua)
		if err != nil {
			return err
		}
		cards.list = nc.list
	}

	return nil
}

func reading(ua *uiautomator.UIAutomator) {
	// 卷动屏幕前稳定1秒
	time.Sleep(time.Second)

	learningTime := RandomLearningTime()
	utils.LearningSwipe(ua, learningTime)
}

func watching() {
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
}

// getCards 文章阅读点击的参照为resourceId，视频也是如此，没有resourceId的文章不去学习
// 根据resourceId获取到所有元素列表，依据元素坐标进行点击操作
func (c *cards) GetCards(ua *uiautomator.UIAutomator) (*cards, error) {
	/*
		// 阅读前确保返回了主页
		err := utils.BackHome(ua)
		if err != nil {
			return nil, err
		}
	*/

	se := uiautomator.Selector{
		"resourceId": "cn.xuexi.android:id/general_card_title_id", // 新闻的resourceId
	}

	element := ua.GetElementBySelector(se)
	count, err := element.Count()
	if err != nil {
		return nil, err
	}

	fmt.Println("waiting for news list...")
	for i := 0; i < count; i++ {
		fmt.Print(".")
		card := element.Eq(i)
		cardText, err := card.GetText()
		if err != nil {
			return nil, err
		}
		cardCenter, err := card.Center(nil)
		if err != nil {
			return nil, err
		}

		c.list = append(c.list, generalCard{cardText, cardCenter})
	}

	return c, nil
}

// CardsSwipe 当初始化一个cards后，使用该方法用来卷动屏幕并返回一个新的cards
func (c *cards) cardSwipe(ua *uiautomator.UIAutomator) (nc *cards, err error) {
	// 根据最后一个card来获取滑动距离
	pStart := uiautomator.Position{
		X: 540,
		Y: 1700,
	}
	pEnd := uiautomator.Position{
		X: pStart.X,
		Y: pStart.Y - 1400,
	}
	ua.Swipe(&pStart, &pEnd, 150)
	time.Sleep(5)

	nc = &cards{}
	nc, err = nc.GetCards(ua)
	if err != nil {
		return nil, err
	}

	return nc, nil
}
