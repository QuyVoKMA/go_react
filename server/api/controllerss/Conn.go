package controllerss

import (
	"fmt"
	"log"
	"net/http"
	"learninggolangweb/gitforum/server/api/controllerss"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// struct là kiểu dữ liệu người dùng tự định nghĩa, nó có tính kế thừa.
// Thằng struc này thay thế cho class
type Server struct {
	DB *gorm.DB
	Router *gin.Engine
}

var errList = make(map[string]string)


//Method trong này là hàm func dưới khai báo riêng cho một kiểu dữ liệu
// dặt biệt này, kiểu dữ liệu này dươc
//
func (server *Server) Initialize(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string){
	var err error
	if Dbdriver == "mysql" {
		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)
		server.DB, err = gorm.Open(Dbdriver, DBURL)
		if err != nil {
			fmt.Printf("Không thể kết nối tới database %s", Dbdriver)
			log.Fatal("Đây là lỗi:", err)
		} else {
			fmt.Printf("Đã kết nối thành công database %s", Dbdriver)
		}
	} else if Dbdriver == "postgres" {
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
		server.DB, err = gorm.Open(Dbdriver, DBURL)
		if err != nil {
			fmt.Printf("Không thể kết nối tới database %s", Dbdriver)
			log.Fatal("Đây là lỗi kết nối tới database postgres :", err)
		} else {
			fmt.Printf("Đã kết nối thành công đến database postgres %s", Dbdriver)
		}
	} else {
		fmt.Println("Không thấy Driver")
	}

	// Chuyển đổi dữ liệu.
	server.DB.Debug().AutoMigrate(
		&models.Users{},
		&models.Product_Clothes{},
		
	)
	server.Router = gin.Default()
	server.Router.Use()

}

func (server *Server)RunPort(addr string){
	log.Fatal(http.ListenAndServe(addr, server.Router))
}