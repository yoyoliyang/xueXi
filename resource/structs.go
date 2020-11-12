// Package resource 包定义了一些常用的资源
package resource

// AppPackageName app包名称，用来检测当前运行app的凭据
var AppPackageName = "cn.xuexi.android"

// Activity 映射对应了不同页面的activity，用作为区分当前页面的凭据
var Activity = map[string]string{
	"home":       "com.alibaba.android.rimet.biz.home.activity.HomeActivity",
	"score":      "com.alibaba.android.rimet.biz.home.activity.MineHomeActivity",
	"learnScore": "com.alibaba.lightapp.runtime.activity.CommonWebViewActivity",
	"news":       "com.alibaba.android.uc.base.navi.window2.Window2Activity",
}
