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
	"xueXi/utils"

	ug "github.com/trazyn/uiautomator-go"
)

func main() {

	if len(os.Args) == 2 {
		ua := ug.New(&ug.Config{
			Host: "192.168.1.52",
			Port: 7912,
		})

		err := utils.BackHome(ua)
		checkErr(err)
		if os.Args[1] == "1" {
			// 阅读和视听学习
			learnList := [...]string{"news", "video"}
			for index, item := range learnList {
				switch item {
				case "news":
					err = titleClick(ua, "综合")
					checkErr(err)
				case "video":
					err = titleClick(ua, "电视台")
					checkErr(err)
					err = titleClick(ua, "联播频道")
					checkErr(err)
				}

				err = learning.Learning(ua, item)
				checkErr(err)

				if index+1 != len(learnList) {
					time.Sleep(time.Second * 2)
				}
			}
		}
		if os.Args[1] == "2" {
			err = learning.AnswerTheQuestion(ua)
		}
	} else {
		fmt.Printf("%v 1 or 2 来学习或回答问题\n", os.Args[0])
	}

}

func checkErr(e error) {
	if e != nil {
		log.Fatal(e)
	}
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
