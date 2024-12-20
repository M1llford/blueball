package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"web2/controllers"
	"web2/dao/mysql"
	"web2/dao/redis"
	"web2/logger"
	"web2/pkg/snowflake"
	"web2/routes"
	"web2/settings"
)

func main() {

	var filePath string
	flag.StringVar(&filePath, "filePath", "./conf/config.json", "路径")
	//解析命令行
	flag.Parse()
	fmt.Println(filePath)
	//返回命令行参数后的其他参数
	//fmt.Println(flag.Args())
	////返回命令行参数后的其它参数个数
	//fmt.Println(flag.NArg())
	////返回使用的命令行参数个数
	//fmt.Println(flag.NFlag())
	//  1.加载配置
	if err := settings.Init(filePath); err != nil {
		fmt.Printf("init settings failed, err:%v\n", err)
		return
	}
	//  2.初始化日志
	if err := logger.Init(settings.Conf.Logconfig); err != nil {
		fmt.Printf("init settings failed, err:%v\n", err)
		return
	}
	defer zap.L().Sync()
	zap.L().Debug("logger init success")
	//  3.初始化MYSQL链接
	if err := mysql.Init(settings.Conf.Mysqlconfig); err != nil {
		fmt.Printf("init settings failed, err:%v\n", err)
		return
	}
	defer mysql.Close()

	//4.初始化redis链接
	if err := redis.Init(settings.Conf.Redisconfig); err != nil {
		fmt.Printf("init settings failed, err:%v\n", err)
		return
	}
	defer redis.Close()
	//雪花ID生成器
	if err := snowflake.Init(settings.Conf.StartTime, settings.Conf.MachineID); err != nil {
		fmt.Printf("StartTime:%s\n", settings.Conf.StartTime)
		fmt.Printf("init snowflake failed, err:%v\n", err)

	}

	//
	if err := controllers.InitTrans("zh"); err != nil {
		zap.L().Error("Init validator failed，err:", zap.Error(err))
		fmt.Printf("Init validator failed, err:%v\n", err)
	}

	//  5.注册路由
	r := routes.Setup(settings.Conf.Gin_mode)

	//  6.启动服务
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.GetInt("app.port")),
		Handler: r,
	}

	go func() {
		// 开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	zap.L().Info("Shutdown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown: ", zap.Error(err))
	}

	zap.L().Info("Server exiting")

}
