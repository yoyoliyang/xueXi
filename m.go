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
	"time"

	"xueXi/learning"
	"xueXi/utils"

	ug "github.com/trazyn/uiautomator-go"
)

func main() {

	ua := ug.New(&ug.Config{
		Host: "192.168.1.52",
		Port: 7912,
	})

	app, _ := ua.GetCurrentApp()
	fmt.Println(app.Package)
	fmt.Println(app.Activity)

	err := utils.BackHome(ua)
	checkErr(err)

	err = utils.ClickPosition(ua, "home")
	checkErr(err)

	// 学习要闻
	err = utils.ClickPosition(ua, "general")
	checkErr(err)

	err = learning.Learning(ua, "news")
	checkErr(err)

	time.Sleep(time.Second * 2)

	// 学习视频
	err = utils.ClickPosition(ua, "tv")
	checkErr(err)

	err = utils.ClickPosition(ua, "tvNews")
	checkErr(err)

	err = learning.Learning(ua, "video")
	checkErr(err)

}

func checkErr(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
