package pkg

import "time"

var (
	debounce = make(map[string]struct{})
)

// CreateDebounce 創建防斗方法
//
//	name 命名規則 [pkg][方法][唯一ID]
//	duration 時間
//	handler 時間到執行
//	return 是否創建成功
func CreateDebounce(name string, duration time.Duration, handler ...func()) bool {
	if _, ok := debounce[name]; ok {
		return false
	}

	debounce[name] = struct{}{}

	go func() {
		timer := time.NewTimer(duration)
		<-timer.C
		timer.Stop()
		delete(debounce, name)
		if len(handler) > 0 {
			for _, f := range handler {
				f()
			}
		}
	}()

	return true
}
