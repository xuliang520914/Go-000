package main

import (
	"container/ring"
	"fmt"
	"math/rand"
	"time"
)

type slidingWindowLimiter struct {
	maxRequest   int           //单位时间最大连接数
	unitTime     time.Duration //单位时间
	windowTime   time.Duration //窗口时间
	windowsCount int           //窗口个数
}

type window struct {
	index              int //窗口编号
	maxRequestCount    int //窗口最大可请求数= 单位时间最大连接数 - 其余全部窗口之前已请求数之和
	handleRequestCount int //窗口当前已处理请求数
}

func NewSlidingWindowLimiter(maxRequest, windowsCount int, unitTime time.Duration) *slidingWindowLimiter {
	return &slidingWindowLimiter{
		maxRequest:   maxRequest,
		unitTime:     unitTime,
		windowTime:   unitTime / time.Duration(windowsCount),
		windowsCount: windowsCount,
	}
}

var requestChan chan string //从中读取请求的url
var request []string        //存储全部接受的请求

func init() {
	rand.Seed(42)
	requestChan = make(chan string, 0)
	request = make([]string, 0)
}

func main() {
	swl := NewSlidingWindowLimiter(100, 10, 10*time.Second)

	//利用循环队列来进行窗口滑动
	r := ring.New(swl.windowsCount)
	for i := 0; i < r.Len(); i++ {
		r.Value = &window{i, swl.maxRequest, 0}
		r = r.Next()
	}
	// 定时器设定时间为窗口时间
	ticker := time.NewTicker(swl.windowTime)
	defer ticker.Stop()

	go getRequest()

	w := r.Value.(*window)
	for {
		select {
		case url := <-requestChan: //不断读取请求
			if w.maxRequestCount > 0 && w.handleRequestCount < w.maxRequestCount {
				request = append(request, url) // 模拟接受请求
				w.handleRequestCount++
				fmt.Printf("this %d window, can handle request max:%d, current has handled request:%d\n", w.index, w.maxRequestCount, w.handleRequestCount)
				continue
			}
		case <-ticker.C:
			fmt.Println("one windowtime has passed ")
			r = r.Move(1) //窗口向后滑动
			w = r.Value.(*window)
			w.maxRequestCount = swl.maxRequest - (getRequestCountSum(r) - w.handleRequestCount) //计算当前窗口能接受请求的最大值
			w.handleRequestCount = 0                                                            // 重置当前窗口已接受请求数量
		}
	}
}

// 模拟收到请求
func getRequest() {
	for {
		time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)
		requestChan <- "url"
	}
}

// 返回全部窗口已处理请求总数
func getRequestCountSum(r *ring.Ring) int {
	var sum int
	for i := 0; i < r.Len(); i++ {
		w := r.Value.(*window)
		sum += w.handleRequestCount
		r = r.Next()
	}
	return sum
}
