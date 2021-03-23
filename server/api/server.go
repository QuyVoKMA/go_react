package api

import (
	"learninggolangweb/gitforum/server/api/controllerss"
	"os"

	"fmt"

	"github.com/joho/godotenv"
)


var server = controllerss.Server{}

// func này dùng để load dữ liệu trong file .env, nó để kiểm tra biến mối trường
func Init() {
	if err := godotenv.Load(); err != nil {
		log.Print("Không thể tìm thấy biến môi trường kết nối database")
	}

}

func Run() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Lỗi! Không thể lấy dữ liệu ENV, %v", err)
	} else {
		fmt.Println("Đã lấy được giá trị")
	}

	// Dữ liệu được đọc và đưa vào func Initialize bên cotrollers.Conn
	server.Initialize(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))

	apiPort := fmt.Sprintf(":%s", os.Getenv("API_PORT"))
	fmt.Printf("Đang lắng nghe cổng %s", apiPort)

	// server.RunPort(apiPort)

}