package bootstrap

import (
	"github.com/qingconglaixueit/wechatbot/handlers"
	"github.com/qingconglaixueit/abing_logger"
	"github.com/eatmoreapple/openwechat"
)

func Run() {
	// 初始化日志
	abing_logger.InitLog(&abing_logger.Config{
		Filename:   "./log/wechatbot.log",
		MaxSize:    200,
		MaxBackups: 10,
		MaxAge:     10,
		Compress:   false,
	})
	//bot := openwechat.DefaultBot()
	bot := openwechat.DefaultBot(openwechat.Desktop) // 桌面模式，上面登录不上的可以尝试切换这种模式

	// 注册消息处理函数
	handler, err := handlers.NewHandler()
	if err != nil {
		abing_logger.SugarLogger.Errorf("register error: %v", err)
		return
	}
	bot.MessageHandler = handler

	// 注册登陆二维码回调
	bot.UUIDCallback = handlers.QrCodeCallBack

	// 创建热存储容器对象
	reloadStorage := openwechat.NewJsonFileHotReloadStorage("storage.json")

	// 执行热登录
	err = bot.HotLogin(reloadStorage, true)
	if err != nil {
		abing_logger.SugarLogger.Warnf("login error: %v ", err)
		return
	}
	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	bot.Block()
}
