package retry

import (
	"fmt"
	"time"
)

/*
Retry 对执行失败的函数，在进行几次重试
*/

// Try 重新尝试函数，如果函数执行失败，则延长时间重试
type Try struct {
	Name     string        //重试任务名称
	F        func() error  //需要重新尝试的函数
	Duration time.Duration //重新尝试的时间间隔
	MaxTimes int           //最大重试次数
}

func NewTry(name string, f func() error, duration time.Duration, maxTimes int) *Try {
	return &Try{
		Name:     name,
		F:        f,
		Duration: duration,
		MaxTimes: maxTimes,
	}
}

// Report 尝试重试的报告
type Report struct {
	Name        string        // 重试任务的名称
	Result      bool          // 函数执行的结果
	Times       int           //重试的次数
	SumDuration time.Duration //总执行时间
	Errs        []error       //函数执行的错误记录
}

func (r *Report) Error() string {
	return fmt.Sprintf("[retry]名称：%s，结果：%v，尝试次数：%v，总时间：%v，错误：%v", r.Name, r.Result, r.Times, r.SumDuration, r.Errs)
}

// Run 尝试重试，返回 chan 可以用于接收尝试报告
// 异步执行重试逻辑，并通过channel返回执行结果
// 返回一个只读channel，可用于接收最终的Report结果
func (try *Try) Run() <-chan Report {
	result := make(chan Report, 1)

	go func() {
		// 确保channel一定会被关闭
		defer close(result)

		// 记录开始时间
		start := time.Now()
		// 存储每次失败的错误信息
		var errs []error

		for i := 0; i < try.MaxTimes; i++ {
			time.Sleep(try.Duration)

			// 执行目标函数
			err := try.F()
			if err == nil {
				// 执行成功，发送成功报告并退出
				result <- Report{
					Name:        try.Name,
					Result:      true,
					Times:       i + 1,
					SumDuration: time.Since(start),
					Errs:        errs,
				}
				return
			}

			// 执行失败，记录错误并继续重试
			errs = append(errs, err)
		}

		// 达到最大重试次数仍失败，发送失败报告
		result <- Report{
			Name:        try.Name,
			Result:      false,
			Times:       try.MaxTimes,
			SumDuration: time.Since(start),
			Errs:        errs,
		}
	}()
	return result
}
