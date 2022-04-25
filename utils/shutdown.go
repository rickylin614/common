package utils

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rickylin614/common/zlog"
)

type GoServer struct {
	srv *http.Server
}

func ServerSet(addr string, handler http.Handler) GoServer {
	return GoServer{srv: &http.Server{
		Addr:    addr,
		Handler: handler,
	}}
}

// 啟動http並附帶優雅關機
func (gs GoServer) Run() {
	srv := gs.srv
	//啟動http listen協程
	go func() {
		zlog.Info("Server start host:", gs.srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zlog.Fatal("main srv listener error", err)
		}
	}()

	// 等待中斷信號來優雅地關閉服務器，為關閉服務器操作設置一個5秒的超時
	quit := make(chan os.Signal, 1) // 創建一個接收信號的通道
	// kill 默認會发送 syscall.SIGTERM 信號 (= -15)
	// kill -2 发送 syscall.SIGINT 信號，我們常用的Ctrl+C就是觸发系統SIGINT信號
	// kill -9 发送 syscall.SIGKILL 信號，但是不能被捕獲，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信號轉发給quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此處不會阻塞
	<-quit                                               // 阻塞在此，當接收到上述兩種信號時才會往下執行
	t := time.Now()                                      // 計算關機總時間使用
	zlog.Info("Begin Shutdown Server...")

	// 創建一個60秒超時的context
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	// 將未處理完的請求處理完再關閉服務），超過N秒就超時退出
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown err: %s , Spend time : %d ms", err, time.Since(t).Milliseconds())
	}

	zlog.Info("Server Exiting")
}
