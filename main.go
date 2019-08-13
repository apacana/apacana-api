package main

import (
	"github.com/apacana/apacana-api/biz/dal/mysql"
	"github.com/apacana/apacana-api/biz/handler"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"os"
	"time"
)

func main() {
	defer func() {
		log.Println("Service stopped.")
	}()
	initLog()
	log.Println("Service starting ......")
	gin.SetMode(gin.ReleaseMode)
	mysql.InitMysql()
	r := gin.Default()
	handler.SetupRouter(r)
	// Listen and Server in 0.0.0.0:8899
	err := r.Run(":8899")
	if err != nil {
		panic(err)
	}
}

func initLog() {
	gin.DisableConsoleColor()
	f, err := os.Create("log/" + time.Now().Format("2006-01-02 15:04:05") + ".log")
	if err != nil {
		panic(err)
	}
	gin.DefaultWriter = io.MultiWriter(f)
}
