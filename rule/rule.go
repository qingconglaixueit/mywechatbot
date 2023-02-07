package rule

import (
	"github.com/qingconglaixueit/abing_logger"
	"io/ioutil"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	STARTTIME = 9
	ENDTIME   = 21
)

type Rule struct{}
// 当前机器人是否工作
var isWork = true
// 当日是否已经次数超限
var degreeOverrun = false
var Grule = &Rule{}
var lock sync.Mutex

func (r *Rule) SetDegreeOverrun(degree bool) {
	lock.Lock()
	defer lock.Unlock()
	degreeOverrun = degree
	return
}
func (r *Rule) GetDegreeOverrun() bool {
	lock.Lock()
	defer lock.Unlock()
	return degreeOverrun
}


func (r *Rule) SetWork(work bool) {
	lock.Lock()
	defer lock.Unlock()
	isWork = work
	return
}
func (r *Rule) GetWork() bool {
	lock.Lock()
	defer lock.Unlock()
	return isWork
}

// 判断时间在今天的早上 9点到 晚上 9 点区间内
func (r *Rule) IsWorkTime(s int, e int) bool {
	if s < 0 || s > 24 {
		s = STARTTIME
	}
	if e < 0 || e > 24 || e <= s {
		e = ENDTIME
	}
	t := time.Now()
	startTime := time.Date(t.Year(), t.Month(), t.Day(), s, 0, 0, 0, time.Local)
	endTime := time.Date(t.Year(), t.Month(), t.Day(), e, 0, 0, 0, time.Local)
	// 判断当前时间是否在当天的 STARTTIME  --  ENDTIME
	if t.Unix() > startTime.Unix() && t.Unix() < endTime.Unix() {
		return true
	}
	return false
}

func (r *Rule) InSlice(str string, sli []string) bool {
	for _, v := range sli {
		if v == str {
			return true
		}
	}
	return false
}

func (r *Rule) WriteNum(num int,numberFile string) error {
	if err := ioutil.WriteFile(numberFile, []byte(strconv.Itoa(num)), 0666); err != nil {
		abing_logger.SugarLogger.Warn("WriteFile error ", err)
		return err
	}
	return nil
}
func(r *Rule)  GetNum(numberFile string) (int, error) {
	data, err := ioutil.ReadFile(numberFile)
	if err != nil {
		abing_logger.SugarLogger.Warn("ReadFile error ", err)
		return 0, err
	}
	tmp := string(data)
	strings.ReplaceAll(tmp,"\n","")
	num ,err :=strconv.Atoi(tmp)
	if err != nil {
		abing_logger.SugarLogger.Warn("Atoi error ", err)
		return 0, err
	}

	return num, nil
}