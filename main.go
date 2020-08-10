package main

import (
	"context"
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/go-programming-tour-book/blog-server/global"
	"github.com/go-programming-tour-book/blog-server/internal/model"
	"github.com/go-programming-tour-book/blog-server/internal/routers"
	"github.com/go-programming-tour-book/blog-server/pkg/logger"
	"github.com/go-programming-tour-book/blog-server/pkg/setting"
	"github.com/go-programming-tour-book/blog-server/pkg/tracer"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	port      string
	runMode   string
	config    string
	isVersion bool
)

// 在go语言中，init方法常用于应用程序内的一些初始化操作，在main方法前自动执行
// 不要滥用init方法，如果init方法过多，则很容易迷失在各个库的init方法中
func init() {
	err := setupSetting()
	if err != nil {
		log.Fatalf("init.setupSetting err: %v", err)
	}

	err = setupDBEngine()
	if err != nil {
		log.Fatalf("init.setupDBEngine err: %v", err)
	}

	err = setupLogger()
	if err != nil {
		log.Fatalf("init.setupLogger err: %v", err)
	}

	err = setupTracer()
	if err != nil {
		log.Fatalf("init setupTracer err: %v", err)
	}
}

func main() {

	/* 测试日志记录可用性 */
	//global.Logger.Infof("%s: go-programming-tour-book/%s", "Ryan", "blog-server")

	gin.SetMode(global.SeverSetting.RunMode)
	router := routers.NewRouter()
	s := &http.Server{
		Addr:           ":" + global.SeverSetting.HttpPort,
		Handler:        router,
		ReadTimeout:    global.SeverSetting.ReadTimeout,
		WriteTimeout:   global.SeverSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	//_ = s.ListenAndServe()

	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("s.ListenAndServer err: %v", err)
		}
	}()
	// 等待信号中断
	quit := make(chan os.Signal)

	// 接收syscall.SIGINT和syscall.SIGTERM信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server")

	// 最大时间控制，通知该服务端他有5s的时间来处理原有的请求
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

func setupSetting() error {
	settings, err := setting.NewSetting()

	if err != nil {
		return err
	}
	err = settings.ReadSection("Server", &global.SeverSetting)
	if err != nil {
		return err
	}

	err = settings.ReadSection("App", &global.AppSetting)
	if err != nil {
		return err
	}

	err = settings.ReadSection("Database", &global.DatabaseSetting)
	if err != nil {
		return err
	}

	err = settings.ReadSection("JWT", &global.JWTSetting)
	if err != nil {
		return err
	}

	err = settings.ReadSection("Email", &global.EmailSetting)
	if err != nil {
		return err
	}

	err = setupFlag()
	if err != nil {
		log.Fatalf("init.setupFlag err: %v", err)
	}

	global.JWTSetting.Expire *= time.Second
	global.SeverSetting.ReadTimeout *= time.Second
	global.SeverSetting.WriteTimeout *= time.Second

	if port != ""{
		global.SeverSetting.HttpPort = port
	}

	if runMode != ""{
		global.SeverSetting.RunMode = runMode
	}

	return nil
}

// 注意global.DBEngine, err = model这里，不是 := （短变量）
// 否则会存在很大的问题，由于 ：= 会重新声明并创建左侧新的局部变量
// 因此在其他包中调用global.DBEngine时，它仍然是nil
func setupDBEngine() error {
	var err error
	global.DBEngine, err = model.NewDBEngine(global.DatabaseSetting)
	if err != nil {
		return err
	}
	return nil
}

func setupLogger() error {
	global.Logger = logger.NewLogger(&lumberjack.Logger{
		Filename:  global.AppSetting.LogSavePath + "/" + global.AppSetting.LogFileName + global.AppSetting.LogFileExt,
		MaxSize:   600,
		MaxAge:    10,
		LocalTime: true,
	}, "", log.LstdFlags).WithCaller(2)
	return nil
}

func setupTracer() error {
	jaegerTracer, _, err := tracer.NewJaegerTracer("blog-server", "127.0.0.1:6831")
	if err != nil {
		return err
	}
	global.Tracer = jaegerTracer
	return nil
}

func setupFlag() error {
	flag.StringVar(&port, "port", "", "启动端口")
	flag.StringVar(&runMode, "mode", "", "启动模式")
	flag.StringVar(&config, "config", "configs/", "指定要使用的配置文件路径")
	flag.BoolVar(&isVersion, "version", false, "编译信息")
	flag.Parse()

	return nil
}
