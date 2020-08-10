package main

import (
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
	"time"
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
	_ = s.ListenAndServe()
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

	global.JWTSetting.Expire *= time.Second
	global.SeverSetting.ReadTimeout *= time.Second
	global.SeverSetting.WriteTimeout *= time.Second

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
