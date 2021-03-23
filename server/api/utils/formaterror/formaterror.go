package formaterror

import (
	"strings"
)

var errorMessages = make(map[string]string)

var err error

func FormatError(errString string) map[string]string {

	if strings.Contains(errString, "username") {
		errorMessages["Taken_username"] = "Tên người dùng đã được sữ dụng"
	}

	if strings.Contains(errString, "email") {
		errorMessages["Taken_email"] = "Email đã được sữ dụng"

	}
	if strings.Contains(errString, "title") {
		errorMessages["Taken_title"] = "Tiêu đề đã được sữ dụng"

	}
	if strings.Contains(errString, "hashedPassword") {
		errorMessages["Incorrect_password"] = "Password sai!!!"
	}
	if strings.Contains(errString, "record not found") {
		errorMessages["No_record"] = "Không tìm thấy bản ghi"
	}

	if strings.Contains(errString, "double like") {
		errorMessages["Double_like"] = "You cannot like this post twice"
	}

	if len(errorMessages) > 0 {
		return errorMessages
	}

	if len(errorMessages) == 0 {
		errorMessages["Incorrect_details"] = "Chi tiết không chính xác"
		return errorMessages
	}

	return nil
}
