package main

import (
	"context"
	"log"
	"net"
	"testing"
	"time"
)

func Test_slice(t *testing.T) {
	s := make([]int, 10, 12)
	s1 := s[8:]
	changeSlice(s1)
	println(&s1)
	t.Logf("s: %v, len of s: %d, cap of s: %d", s, len(s), cap(s))
	t.Logf("s1: %v, len of s1: %d, cap of s1: %d", s1, len(s1), cap(s1))
}

func changeSlice(s1 []int) {
	println(&s1)
	s1 = append(s1, 10)
	println(&s1)
}

func MonthScope(stamp int64) (int64, int64) {
	now := time.Unix(stamp, 0)
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	return firstOfMonth.Unix(), lastOfMonth.Unix() + 3600*24 - 1
}

func Test_cancelctx(t *testing.T) {
	ctx := context.Background()
	cancel, cancelFunc := context.WithCancel(ctx)
	go func() {
		time.Sleep(time.Second)
		cancelFunc()
	}()

	for {
		select {
		case <-cancel.Done():
			log.Println(cancel.Err())
			return
		default:
			log.Println("waiting")
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func Test_go_ticker(t *testing.T) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	ticker := time.NewTicker(3 * time.Second)
	go func() {
		time.Sleep(10 * time.Second)
		cancelFunc()
	}()
	for  {
		select {
		case <-ctx.Done():
			println("stop")
			return
		case <-ticker.C:
			println("tiktok")
		default:
		}
	}
}

func Test_net(t *testing.T) {

	println("test...")
}




func serve(conn net.Conn){
	defer conn.Close()
	var buf []byte
	// 读取连接中的数据
	_,_ = conn.Read(buf)
	// ...
}
