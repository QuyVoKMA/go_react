package models

import (
	"errors"
	"html"
	"log"
	"os"
	"learninggolangweb/gitforum/server/api/security"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
)

type Users struct {
	ID uint32 	`gorm:"primary_key;auto_increment" json:"id"`
	Username string `gorm:"size:255; not null; unique" json:"username"`
	Email string	`gorm:"size:255; not null; unique" json:"email"`
	Password string	`gorm:"size:255;not null; unique" json:"password"`
	AvartarPath string `gorm:"size:255;null"json:"avartar_path"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdateAt time.Time	`gorm:"default:CURRENT_TIMESTAMP" json:"update_at"`
}

// Trước khi lưu thì hash password
func (u *Users) BeforeSave() error {
	hashedPassword, err := security.Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}


// Loại bở các lỗi, các lỗ hỗng về ký tự trước khi lưu
func (u *Users) Prepare() {
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.CreatedAt = time.Now()
	u.UpdateAt = time.Now()
}


func (u *Users) AfterFind() (err error) {
	if err != nil {
		return err
	}
	if u.AvartarPath != "" {
		u.AvartarPath = os.Getenv("DO_SPACES_URL") + u.AvartarPath
	}
	//dont return the user password
	// u.Password = ""
	return nil
}

// Xác nhận
func (u *Users) Validate(action string) map[string]string {
	var errorMessages = make(map[string]string)
	var err error

	switch strings.ToLower(action) {
	case "update":
		if u.Email == "" {
			err = errors.New("Bắt buộc nhập Email")
			errorMessages["Required_email"] = err.Error()
		}
		if u.Email != "" {
			if err = checkmail.ValidateFormat(u.Email); err != nil {
				err = errors.New("Email sai")
				errorMessages["Invalid_email"] = err.Error()
			}
		}

	case "login":
		if u.Password == "" {
			err = errors.New("Bắt buộc nhập Password")
			errorMessages["Required_password"] = err.Error()
		}
		if u.Email == "" {
			err = errors.New("Bắt buộc nhập Email")
			errorMessages["Required_email"] = err.Error()
		}
		if u.Email != "" {
			if err = checkmail.ValidateFormat(u.Email); err != nil {
				err = errors.New("Email Sai")
				errorMessages["Invalid_email"] = err.Error()
			}
		}
	case "forgotpassword":
		if u.Email == "" {
			err = errors.New("Bắt buộc nhập Email")
			errorMessages["Required_email"] = err.Error()
		}
		if u.Email != "" {
			if err = checkmail.ValidateFormat(u.Email); err != nil {
				err = errors.New("Email Sai")
				errorMessages["Invalid_email"] = err.Error()
			}
		}
	default:
		if u.Username == "" {
			err = errors.New("Bắt buộc nhập Username")
			errorMessages["Required_username"] = err.Error()
		}
		if u.Password == "" {
			err = errors.New("Bắt buộc nhập Password")
			errorMessages["Required_password"] = err.Error()
		}
		if u.Password != "" && len(u.Password) < 6 {
			err = errors.New("Mật khẩu bắc buộc phải có ít nhất 6 ký tự")
			errorMessages["Invalid_password"] = err.Error()
		}
		if u.Email == "" {
			err = errors.New("Bắt buộc nhập Email")
			errorMessages["Required_email"] = err.Error()

		}
		if u.Email != "" {
			if err = checkmail.ValidateFormat(u.Email); err != nil {
				err = errors.New("Email Sai")
				errorMessages["Invalid_email"] = err.Error()
			}
		}
	}
	return errorMessages
}

func (u *Users) SaveUser(db *gorm.DB) (*Users, error) {

	var err error
	err = db.Debug().Create(&u).Error
	if err != nil {
		return &Users{}, err
	}
	return u, nil
}

// THE ONLY PERSON THAT NEED TO DO THIS IS THE ADMIN, SO I HAVE COMMENTED THE ROUTES, SO SOMEONE ELSE DONT VIEW THIS DETAILS.
// Chỉ có Admin mới làm được điều này
func (u *Users) FindAllUsers(db *gorm.DB) (*[]Users, error) {
	var err error
	users := []Users{}
	err = db.Debug().Model(&Users{}).Limit(100).Find(&users).Error
	if err != nil {
		return &[]Users{}, err
	}
	return &users, err
}

func (u *Users) FindUserByID(db *gorm.DB, uid uint32) (*Users, error) {
	var err error
	err = db.Debug().Model(Users{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &Users{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Users{}, errors.New("Không tìm thấy User")
	}
	return u, err
}

func (u *Users) UpdateAUser(db *gorm.DB, uid uint32) (*Users, error) {

	if u.Password != "" {
		// To hash the password
		err := u.BeforeSave()
		if err != nil {
			log.Fatal(err)
		}

		db = db.Debug().Model(&Users{}).Where("id = ?", uid).Take(&Users{}).UpdateColumns(
			map[string]interface{}{
				"password":  u.Password,
				"email":     u.Email,
				"update_at": time.Now(),
			},
		)
	}

	db = db.Debug().Model(&Users{}).Where("id = ?", uid).Take(&Users{}).UpdateColumns(
		map[string]interface{}{
			"email":     u.Email,
			"update_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &Users{}, db.Error
	}

	// This is the display the updated user
	err := db.Debug().Model(&Users{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &Users{}, err
	}
	return u, nil
}

func (u *Users) UpdateAUserAvatar(db *gorm.DB, uid uint32) (*Users, error) {
	db = db.Debug().Model(&Users{}).Where("id = ?", uid).Take(&Users{}).UpdateColumns(
		map[string]interface{}{
			"avatar_path": u.AvartarPath,
			"update_at":   time.Now(),
		},
	)
	if db.Error != nil {
		return &Users{}, db.Error
	}
	// This is the display the updated user
	err := db.Debug().Model(&Users{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &Users{}, err
	}
	return u, nil
}

func (u *Users) DeleteAUser(db *gorm.DB, uid uint32) (int64, error) {

	db = db.Debug().Model(&Users{}).Where("id = ?", uid).Take(&Users{}).Delete(&Users{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (u *Users) UpdatePassword(db *gorm.DB) error {

	// To hash the password
	err := u.BeforeSave()
	if err != nil {
		log.Fatal(err)
	}

	db = db.Debug().Model(&Users{}).Where("email = ?", u.Email).Take(&Users{}).UpdateColumns(
		map[string]interface{}{
			"password":  u.Password,
			"update_at": time.Now(),
		},
	)
	if db.Error != nil {
		return db.Error
	}
	return nil
}