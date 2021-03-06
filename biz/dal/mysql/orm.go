package mysql

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"os"
)

var (
	DB    *gorm.DB
	DBErr error
)

func InitMysql() {
	DB, DBErr = gorm.Open("mysql", "root:@tcp(127.0.0.1:3306)/apacana")
	if DBErr != nil {
		log.Println("connect database fail:", DBErr)
		os.Exit(1)
		return
	}
	log.Println("connect database success")
}

func Insert(tx *gorm.DB, data interface{}) error {
	if tx == nil {
		tx = DB
	}

	if err := tx.Create(data).Error; err != nil {
		_, _ = fmt.Fprintf(gin.DefaultWriter, "insert mysql error: %s\n", err)
		return err
	}
	return nil
}
