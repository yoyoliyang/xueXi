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
var learningTime = 65

func Test() {
	rand.Seed(42)
	fmt.Println(rand.Intn(learningTime))
}

func Reading(ua *uiautomator.UIAutomator) error {
	log.Println("starting reading:")

	cards := &cards{}
	cards, err := cards.GetCards(ua)
	if err != nil {
		return err
	}

	for index, card := range cards.list {
		fmt.Println(index, card.title)
		// 点击新闻标题进入
		err := ua.Click(card.position)
		if err != nil {
			return err
		}

		// 卷动屏幕
		utils.LearningSwap(ua, learningTime)

		// 返回首页
		err = utils.BackHome(ua)
		if err != nil {
			return err
		}
	}

	return nil
}

// Watching 视频学习频道为"联播频道"，和news阅读一样，依据resourceId获取cards
func Watching(ua *uiautomator.UIAutomator) error {
	log.Println("starting watching:")

	cards := &cards{}
	cards, err := cards.GetCards(ua)
	if err != nil {
		return err
	}

	for index, card := range cards.list {
		fmt.Println(index, card.title)
		// 点击标题进入
		err := ua.Click(card.position)
		if err != nil {
			return err
		}

		go func() {
			for i := 0; i < learningTime; i++ {
				fmt.Print(".")
				time.Sleep(time.Second)
			}
		}()

		time.Sleep(time.Duration(learningTime) * time.Second)

		// 返回首页
		err = utils.BackHome(ua)
		if err != nil {
			return err
		}
	}
	return nil
}

// getCards 文章阅读点击的参照为resourceId，视频也是如此，没有resourceId的文章不去学习
// 根据resourceId获取到所有元素列表，依据元素坐标进行点击操作
func (c *cards) GetCards(ua *uiautomator.UIAutomator) (*cards, error) {
	// 阅读前确保返回了主页
	err := utils.BackHome(ua)
	if err != nil {
		return nil, err
	}

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
