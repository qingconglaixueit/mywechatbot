// @Author Bing 
// @Date 2023/2/6 17:33:00 
// @Desc
package handlers

import (
	"github.com/qingconglaixueit/abing_logger"
	"github.com/qingconglaixueit/wechatbot/config"
	"github.com/qingconglaixueit/wechatbot/model/wordfilter"
	"github.com/qingconglaixueit/wechatbot/rule"
)

// 记录请求 gpt 次数
func processNumberFile() {
	if ToTalNumber != 0 {
		ToTalNumber = ToTalNumber + 1
	} else {
		tmp, err := rule.Grule.GetNum(numberFile)
		if err != nil {
			abing_logger.SugarLogger.Warnf("rule.Grule.GetNum error ", err)
		}
		ToTalNumber = tmp + 1
	}
	if err := rule.Grule.WriteNum(ToTalNumber, numberFile); err != nil {
		abing_logger.SugarLogger.Warnf("rule.Grule.WriteNum error ", err)
	}
}

// 记录当日请求次数
func processCurrentReqTimes() {
	if CurrentDayDegree >= config.LoadConfig().CurrentMaxReq {
		rule.Grule.SetDegreeOverrun(true)
		abing_logger.SugarLogger.Warnf("cureent max req  :%d>=%d", CurrentDayDegree, config.LoadConfig().CurrentMaxReq)
		return
	}

	if CurrentDayDegree != 0 {
		CurrentDayDegree = CurrentDayDegree + 1
	} else {
		tmp, err := rule.Grule.GetNum(tmpReqFile)
		if err != nil {
			abing_logger.SugarLogger.Warnf("rule.Grule.GetNum error ", err)
		}
		CurrentDayDegree = tmp + 1
	}

	if err := rule.Grule.WriteNum(CurrentDayDegree, tmpReqFile); err != nil {
		abing_logger.SugarLogger.Warnf("rule.Grule.WriteNum error ", err)
	}
}

// 重置 每日请求 gpt 的次数
func ReSetCurrentReqTimes() {
	if err := rule.Grule.WriteNum(0, tmpReqFile); err != nil {
		abing_logger.SugarLogger.Warnf("rule.Grule.WriteNum error ", err)
		return
	}
	// 重置全局变量
	rule.Grule.SetDegreeOverrun(false)
	abing_logger.SugarLogger.Info("ReSetCurrentReqTimes successfully ... ")
}

// 校验是否有敏感词
func IsWordFilter(str string) bool {
	res := wordfilter.Filter.FindAll(str)
	abing_logger.SugarLogger.Info("wordfilter : ", res)
	if len(res) > 0 {
		return true
	}
	return false
}
