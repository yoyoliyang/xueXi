// 学习强国自动学习脚本
// 2020-11-10

/*
流程：
1、根据学习积分来获取当前需要学习的项目
2、根据学习项目进入具体页面
	每个学习项目需要制定单独的学习模块
3、进行学习
*/

package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"xueXi/learning"
	"xueXi/notice"
	"xueXi/resource"
	"xueXi/utils"

	ug "github.com/trazyn/uiautomator-go"
)

func main() {

	if len(os.Args) == 2 {
		ua := ug.New(&ug.Config{
			Host: "192.168.1.52",
			Port: 7912,
		})

		// 唤醒屏幕
		err := ua.WakeUp()
		if err != nil {
			log.Fatal(err)
		}

		// 检测当前是否为app界面
		appInfo, err := ua.GetCurrentApp()
		if err != nil {
			log.Fatal(err)
		}
		if appInfo.Package != resource.AppPackageName {
			fmt.Println("启动学习强国app")
			err := ua.AppStart(resource.AppPackageName)
			if err != nil {
				log.Fatal(err)
			}
			time.Sleep(time.Second * 7)
		}

		// 登录操作
		err = utils.Login(ua)
		if err != nil {
			log.Fatal(err)
		}

		// home界面刷新
		err = utils.BackHome(ua)
		if err != nil {
			log.Fatal(err)
		}

		// 获取当前积分
		score, err := utils.GetCurrentScore(ua)
		fmt.Println("当前学习积分: ", score)

		switch os.Args[1] {

		case "1":
			// 阅读和视听学习
			fmt.Println("开始学习新闻")
			err = titleClick(ua, "综合")
			if err != nil {
				log.Fatal(err)
			}

			err = learning.Learning(ua, "news")
			if err != nil {
				log.Fatal(err)
			}

		case "2":
			fmt.Println("开始学习视频")
			err = titleClick(ua, "电视台")
			if err != nil {
				log.Fatal(err)
			}
			err = titleClick(ua, "联播频道")
			if err != nil {
				log.Fatal(err)
			}
			err = learning.Learning(ua, "video")
			if err != nil {
				log.Fatal(err)
			}
		case "3":
			fmt.Println("开始每日答题")
			err = learning.AnswerTheQuestion(ua)
			if err != nil {
				log.Fatal(err)
			}
		}
		// 获取今日学习积分
		if _score, err := utils.GetCurrentScore(ua); err == nil {
			score += _score
			scoreMsg := fmt.Sprintf("今日学习积分: %v", score)
			fmt.Println(scoreMsg)
			notice.IftttNotice(scoreMsg)
		} else {
			log.Fatal(err)
		}
	}
	fmt.Printf("%v 1/2/3来学习或回答问题\n", os.Args[0])

}

func titleClick(ua *ug.UIAutomator, name string) error {
	defer time.Sleep(time.Second * 2)
	position, err := utils.GetSelectorPostion(ua, &ug.Selector{
		"className": "android.widget.TextView",
		"text":      name,
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
