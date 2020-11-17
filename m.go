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

	learnList := [...]string{"news", "video"}
	for index, item := range learnList {
		switch item {
		case "news":
			err = utils.ClickPosition(ua, "综合")
			checkErr(err)
		case "video":
			err = utils.ClickPosition(ua, "电视台")
			checkErr(err)
		}

		err = learning.Learning(ua, item)
		checkErr(err)

		if index+1 != len(learnList) {
			time.Sleep(time.Second * 2)
		}
	}

}

func checkErr(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
