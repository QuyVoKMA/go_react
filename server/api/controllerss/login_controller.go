package controllerss

import (
	
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)



func (server *Server) Login(c *gin.Context) {
	// Xóa lỗi trước khi nó có
	errList = map[string]string{}

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		// Thông báo trạng thái ko thẻ thực thi
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":      http.StatusUnprocessableEntity,
			"first error": "Không thể nhận được yêu cầu",
		})
		return

	}

	user := models.Users{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  "Không thể kiểm soát cơ chế",
		})
		return
	}
	user.Prepare()
	errorMessages := user.Validate("login")
	if len(errorMessages) > 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errorMessages,
		})
		return
	}

	userData, err := server.SignIn(user.Email, user.Password)
	if err != nil {
		formatedError := formaterror.FormatError(err.Error())

		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"Error":  formatedError,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": userData,
	})
}

// Đăng nhập
func (server *Server) SignIn(email, password string) (map[string]interface{}, error) {

	var err error

	userData := make(map[string]interface{})

	user := models.Users{}

	err = server.DB.Debug().Model(models.Users{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		fmt.Println("Đây là lỗi nhận được từ người dùng: ", err)
		return nil, err
	}
	err = security.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		fmt.Println("Đây là lỗi khi băm mật khẩu ", err)
		return nil, err
	}
	token, err := authen.CreateToken(user.ID)
	if err != nil {
		fmt.Println("Đây là lỗi khi tạo mã  token: ", err)
		return nil, err
	}
	userData["token"] = token
	userData["id"] = user.ID
	userData["email"] = user.Email
	userData["avatar_path"] = user.AvartarPath
	userData["username"] = user.Username

	return userData, nil
}
