package security

import "golang.org/x/crypto/bcrypt"



//Password được chuyển vào hàm này và tiến hành hash 
func Hash(password string) ([]byte, error) {

	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// kiếm chứng password có đúng ko bằng cách so sánh
func VerifyPassword(hashedPassword, password string) error {

	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

}