package handlers

import (
	"errors"
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/qingconglaixueit/abing_logger"
	"github.com/qingconglaixueit/wechatbot/config"
	"github.com/qingconglaixueit/wechatbot/gpt"
	"github.com/qingconglaixueit/wechatbot/rule"
	"github.com/qingconglaixueit/wechatbot/service"
	"strings"
)

const (
	HaveARest = "该休息了"
	WorkStr   = "起来嗨"
	WorkTime  = "你的工作时间"
	ResDegree = "lastdegree"

	replyForRest = "好的，我开始睡美容觉了！"
	replyForWork = "主人，我现在活力满满！"

	replyForRest2 = "我已经休息了，你也快睡觉吧！"

	replyPersonal = "不要单独聊我哦，可以到群里互动哦！！"

	replyForRestMaxReq = "今日请求超限了哦！！我回复不动了！！"

	replyWorkTime = "我的工作时间是 %d 点 -- %d 点"

	replyDegree = "今日 gpt 请求剩余次数为：%d，总共 %d"

	replyGptErrStr = "有点蒙了，待我整理一下思绪，你轻点..."

	replyWordFilter = "差点把我整懵了，请使用正确的词汇..."

	replyServerBusy = "服务器忙，请稍后..."

	numberFile = "./number.txt"
	// 每天仅限请求 100 次
	tmpReqFile  = "./tmpReqFile.txt"
	reqMaxTimes = 100
)

// 服务启动以来总共请求 gpt 的次数
var ToTalNumber = 0

// 当日已经请求的次数
var CurrentDayDegree = 0

// VIP 账号
var VipUserList = []string{"Anonymous", "LOqw789"}
var _ MessageHandlerInterface = (*GroupMessageHandler)(nil)

// GroupMessageHandler 群消息处理
type GroupMessageHandler struct {
	// 获取自己
	self *openwechat.Self
	// 群
	group *openwechat.Group
	// 接收到消息
	msg *openwechat.Message
	// 发送的用户
	sender *openwechat.User
	// 实现的用户业务
	service service.UserServiceInterface
}

func GroupMessageContextHandler() func(ctx *openwechat.MessageContext) {
	return func(ctx *openwechat.MessageContext) {
		msg := ctx.Message
		// 获取用户消息处理器
		handler, err := NewGroupMessageHandler(msg)
		if err != nil {
			abing_logger.SugarLogger.Warn(fmt.Sprintf("init group message handler error: %s", err))
			return
		}

		// 处理用户消息
		err = handler.handle()
		if err != nil {
			abing_logger.SugarLogger.Warn(fmt.Sprintf("handle group message error: %s", err))
		}
	}
}

// NewGroupMessageHandler 创建群消息处理器
func NewGroupMessageHandler(msg *openwechat.Message) (MessageHandlerInterface, error) {
	sender, err := msg.Sender()
	if err != nil {
		return nil, err
	}
	group := &openwechat.Group{User: sender}
	groupSender, err := msg.SenderInGroup()
	if err != nil {
		return nil, err
	}

	userService := service.NewUserService(c, groupSender)
	handler := &GroupMessageHandler{
		self:    sender.Self,
		msg:     msg,
		group:   group,
		sender:  groupSender,
		service: userService,
	}
	return handler, nil

}

// handle 处理消息
func (g *GroupMessageHandler) handle() error {
	if g.msg.IsText() {
		return g.ReplyText()
	}
	return nil
}

// ReplyText 发息送文本消到群
func (g *GroupMessageHandler) ReplyText() error {
	abing_logger.SugarLogger.Info(fmt.Sprintf("Received Group %v Text Msg : %v", g.group.NickName, g.msg.Content))
	var (
		err   error
		reply string
	)

	// 1.不是@的不处理
	if !g.msg.IsAt() {
		return nil
	}

	// 2.获取请求的文本，如果为空字符串不处理
	requestText := g.getRequestText()
	if requestText == "" {
		abing_logger.SugarLogger.Warn("user message is null")
		return nil
	}
	// 检查敏感词
	if IsWordFilter(requestText) {
		reply = replyWordFilter
	}else{
		// 获取特殊回复消息
		reply = g.getSpecialReply(requestText)
	}

	if rule.Grule.GetWork() && reply == "" {
		isSvr := true
		if !rule.Grule.IsWorkTime(config.LoadConfig().StartTime, config.LoadConfig().EndTime) {
			isSvr = false
			reply = replyForRest2
		}
		// 非工作时间，仍然服务 vip 用户
		if !rule.Grule.IsWorkTime(config.LoadConfig().StartTime, config.LoadConfig().EndTime) &&
			rule.Grule.InSlice(g.sender.NickName, VipUserList) {
			isSvr = true
			reply = ""
		}
		// 能服务的时候，才请求 chatgpt
		if isSvr {
			// 3.请求GPT获取回复
			// 记录请求 gpt 次数
			processNumberFile()
			// 将当日次数写入到临时文件中，且全局变量 +1
			processCurrentReqTimes()
			// 如果超限
			abing_logger.SugarLogger.Info(" rule.Grule.GetDegreeOverrun() == ", rule.Grule.GetDegreeOverrun())
			if rule.Grule.GetDegreeOverrun() {
				abing_logger.SugarLogger.Info(" rule.Grule.GetDegreeOverrun() == ", rule.Grule.GetDegreeOverrun())
				reply = replyForRestMaxReq
			} else {
				reply, err = gpt.Completions(requestText)
				if err != nil {
					// 2.1 将GPT请求失败信息输出给用户，省得整天来问又不知道日志在哪里。
					errMsg := fmt.Sprintf("gpt request error: %v", err)
					abing_logger.SugarLogger.Info(errMsg)
					_, err = g.msg.ReplyText(replyGptErrStr)
					if err != nil {
						return errors.New(fmt.Sprintf("response group error: %v ", err))
					}
					return err
				}
			}

		}
	} else {
		if reply == "" {
			reply = replyForRest2
		}
	}

	if reply == "" {
		reply = replyGptErrStr
	}

	// 4.设置上下文，并响应信息给用户
	g.service.SetUserSessionContext(requestText, reply)
	_, err = g.msg.ReplyText(g.buildReplyText(reply))
	if err != nil {
		return errors.New(fmt.Sprintf("response user error: %v ", err))
	}

	// 5.返回错误信息
	return err
}

// getRequestText 获取请求接口的文本，要做一些清洗
func (g *GroupMessageHandler) getRequestText() string {
	// 1.去除空格以及换行
	requestText := strings.TrimSpace(g.msg.Content)
	requestText = strings.Trim(g.msg.Content, "\n")

	// 2.替换掉当前用户名称
	replaceText := "@" + g.self.NickName
	requestText = strings.TrimSpace(strings.ReplaceAll(g.msg.Content, replaceText, ""))
	if requestText == "" {
		return ""
	}
	abing_logger.SugarLogger.Info("2222 requestText == ", requestText)
	if requestText == WorkTime || requestText == HaveARest || requestText == WorkStr || requestText == ResDegree {
		return requestText
	}

	// 3.获取上下文，拼接在一起，如果字符长度超出4000，截取为4000。（GPT按字符长度算），达芬奇3最大为4068，也许后续为了适应要动态进行判断。
	sessionText := g.service.GetUserSessionContext()
	if sessionText != "" {
		requestText = sessionText + "\n" + requestText
	}
	if len(requestText) >= 4000 {
		requestText = requestText[:4000]
	}

	// 4.检查用户发送文本是否包含结束标点符号
	punctuation := ",.;!?，。！？、…"
	runeRequestText := []rune(requestText)
	lastChar := string(runeRequestText[len(runeRequestText)-1:])
	if strings.Index(punctuation, lastChar) < 0 {
		requestText = requestText + "？" // 判断最后字符是否加了标点，没有的话加上句号，避免openai自动补齐引起混乱。
	}

	// 5.返回请求文本
	return requestText
}

// buildReply 构建回复文本
func (g *GroupMessageHandler) buildReplyText(reply string) string {
	// 1.获取@我的用户
	atText := "@" + g.sender.NickName
	textSplit := strings.Split(reply, "\n\n")
	if len(textSplit) > 1 {
		trimText := textSplit[0]
		reply = strings.Trim(reply, trimText)
	}
	reply = strings.TrimSpace(reply)
	if reply == "" {
		return atText + " 请求得不到任何有意义的回复，请具体提出问题。"
	}

	// 2.拼接回复,@我的用户，问题，回复
	replaceText := "@" + g.self.NickName
	question := strings.TrimSpace(strings.ReplaceAll(g.msg.Content, replaceText, ""))
	reply = atText + "\n" + question + "\n --------------------------------\n" + reply
	reply = strings.Trim(reply, "\n")

	// 3.返回回复的内容
	return reply
}
// 获取特殊回复消息
func (g *GroupMessageHandler) getSpecialReply(requestText string) (reply string) {
	abing_logger.SugarLogger.Info("requestText == ", requestText)

	if requestText == WorkTime {
		reply = fmt.Sprintf(replyWorkTime, config.LoadConfig().StartTime, config.LoadConfig().EndTime)
	}
	abing_logger.SugarLogger.Info("--------------------", rule.Grule.InSlice(g.sender.NickName, VipUserList))

	// 识别到是 Anonymous 发送过来的消息，且是 “该休息了！”，那么则将全局变量设置为 false，休息状态
	if rule.Grule.InSlice(g.sender.NickName, VipUserList) && requestText == HaveARest {
		abing_logger.SugarLogger.Info("have a rest !!!!!")
		rule.Grule.SetWork(false)
		reply = replyForRest
	}

	if rule.Grule.InSlice(g.sender.NickName, VipUserList) && requestText == WorkStr {
		abing_logger.SugarLogger.Info("work start !!!!!")
		rule.Grule.SetWork(true)
		reply = replyForWork
	}

	abing_logger.SugarLogger.Info("requestText == ", requestText)
	if requestText == ResDegree {
		degree, dErr := rule.Grule.GetNum(tmpReqFile)
		if dErr != nil {
			abing_logger.SugarLogger.Errorf("GetNum error:%s", dErr)
			reply = replyServerBusy
		} else {
			num := config.LoadConfig().CurrentMaxReq - degree
			if num >= 0 {
				abing_logger.SugarLogger.Infof("now res degree :%d", num)
				reply = fmt.Sprintf(replyDegree, num, config.LoadConfig().CurrentMaxReq)
			} else {
				reply = fmt.Sprintf(replyDegree, 0, config.LoadConfig().CurrentMaxReq)
			}
		}
	}
	return
}
