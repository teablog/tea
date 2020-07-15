package initialize

import (
	"bytes"
	"context"
	"dyc/internal/config"
	"dyc/internal/controllers"
	"dyc/internal/db"
	"dyc/internal/derror"
	"dyc/internal/logger"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"strings"
	"time"
)

func Server() {
	var (
		engine *gin.Engine
		err    error
	)
	// 日志
	logger.NewLogger(config.GetLogFD())
	// 数据库
	db.NewElasticsearch(config.GetKey("elasticsearch::address").Strings(","), config.GetKey("elasticsearch::user").String(), config.GetKey("elasticsearch::password").String())
	// 启动gin
	switch config.GetKey("global::env").String() {
	case "prod":
		gin.SetMode(gin.ReleaseMode)
		engine = gin.New()
		engine.Use(recoverWithWrite(config.GetLogFD()))
	default:
		engine = gin.Default()
	}

	// 路由
	controllers.Init(engine)
	controllers.NewRouter(engine)
	port := config.GetKey("server::port").String()
	server := http.Server{
		Addr:     port,
		Handler:  engine,
		ErrorLog: nil,
	}
	go func() {
		logger.Infof("start server 127.0.0.1%s", port)
		if err = server.ListenAndServe(); err != nil {
			logger.Info("web server closed")
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// 关闭服务器
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("web server close failed: %s", err)
	}
	time.Sleep(50 * time.Millisecond)
}

func recoverWithWrite(out io.Writer) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				headers := strings.Split(string(httpRequest), "\r\n")
				for idx, header := range headers {
					current := strings.Split(header, ":")
					if current[0] == "Authorization" {
						headers[idx] = current[0] + ": *"
					}
				}
				if brokenPipe {
					logger.Debugf("%s\n%s%s", err, string(httpRequest), logger.Reset)
					c.Error(err.(error)) // nolint: errcheck
				} else {
					switch err.(type) {
					case derror.Unauthorized:
						logger.Debugf("cookie 验证失败")
						c.JSON(http.StatusOK, gin.H{"msg": "unauthorized", "code": http.StatusUnauthorized})
					default:
						buf := new(bytes.Buffer) // the returned data
						err := errors.WithStack(err.(error))
						fmt.Fprintf(buf, "%+v", err)
						logger.Errorf("[Recovery] panic recovered:\n%s\n%s", strings.Join(headers, "\r\n"), buf.String())
						c.JSON(http.StatusInternalServerError, gin.H{"msg": "服务器出错了", "code": http.StatusInternalServerError})
					}
				}
				c.Abort()
			}
		}()
		c.Next()
	}
}
