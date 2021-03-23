package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)


type Product_Clothes struct {
	ID       uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Title    string    `gorm:"size:255; not null; unique" json:"title"`
	Content  string    `gorm:"text; not null" json:"content"`
	ImageClothes string `gorm:"size:255; not null" json:"image_clothes"`
	Price   float64		`gorm:"not null" json:"price"`
	Author   Users      `json:"author"`
	AuthorID uint32    `gorm:"not nul" json:"author_id"`
	CreateAt time.Time     `gorm:"default:CURRENT_TIMESTAMP" json:"create_at"`
	UpdateAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"update_at"`
}

func (pc *Product_Clothes) Prepare(){
	pc.Title = html.EscapeString(strings.TrimSpace(pc.Title))
	pc.Content = html.EscapeString(strings.TrimSpace(pc.Content))
	pc.Author = Users{}
	pc.CreateAt = time.Now()
	pc.UpdateAt = time.Now()
}

// Xác thực
func (p *Product_Clothes) Validate() map[string]string {

	var err error

	var errorMessages = make(map[string]string)

	if p.Title == "" {
		err = errors.New("Bắt buộc phải có tiêu đề")
		errorMessages["Required_title"] = err.Error()

	}
	if p.Content == "" {
		err = errors.New("Bắt buộc phải có nội dung")
		errorMessages["Required_content"] = err.Error()

	}
	if p.AuthorID < 1 {
		err = errors.New("Bắt buộc phải có tác giả")
		errorMessages["Required_author"] = err.Error()
	}
	return errorMessages
}

// Lưu sản phẩm.
func (p *Product_Clothes) SaveProductClothes(db *gorm.DB) (*Product_Clothes, error) {
	var err error
	// Lệnh này dùng để tạo một danh mục sản phẩm
	err = db.Debug().Model(&Product_Clothes{}).Create(&p).Error
	if err != nil {
		return &Product_Clothes{}, err
	}
	//Lúc này tạo tự tạo một id tác giả khi sp được tạo
	if p.ID != 0 {
		err = db.Debug().Model(&Users{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Product_Clothes{}, err
		}
	}
	return p, nil
}

// Tìm tất cả các post trong 1 lần
func (p *Product_Clothes) FindAllProductClothes(db *gorm.DB) (*[]Product_Clothes, error) {
	var err error
	product_Clothes := []Product_Clothes{}
	err = db.Debug().Model(&Product_Clothes{}).Limit(100).Order("created_at desc").Find(&product_Clothes).Error
	if err != nil {
		return &[]Product_Clothes{}, err
	}
	if len(product_Clothes) > 0 {
		for i, _ := range product_Clothes {
			err := db.Debug().Model(&Users{}).Where("id = ?", product_Clothes[i].AuthorID).Take(&product_Clothes[i].Author).Error
			if err != nil {
				return &[]Product_Clothes{}, err
			}
		}
	}
	return &product_Clothes, nil
}

// Tìm theo ID
func (p *Product_Clothes) FindProductClothesByID(db *gorm.DB, pid uint64) (*Product_Clothes, error) {
	var err error
	err = db.Debug().Model(&Product_Clothes{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &Product_Clothes{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&Users{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Product_Clothes{}, err
		}
	}
	return p, nil
}

// Cập nhật 
func (p *Product_Clothes) UpdateAProductClothes(db *gorm.DB) (*Product_Clothes, error) {

	var err error

	err = db.Debug().Model(&Product_Clothes{}).Where("id = ?", p.ID).Updates(Product_Clothes{Title: p.Title, Content: p.Content, UpdateAt: time.Now()}).Error
	if err != nil {
		return &Product_Clothes{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&Users{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Product_Clothes{}, err
		}
	}
	return p, nil
}

//Xóa
func (p *Product_Clothes) DeleteAProductClothes(db *gorm.DB) (int64, error) {

	db = db.Debug().Model(&Product_Clothes{}).Where("id = ?", p.ID).Take(&Product_Clothes{}).Delete(&Product_Clothes{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

// Tìm tác giả của đang.
func (p *Product_Clothes) FindUserProductClothes(db *gorm.DB, uid uint32) (*[]Product_Clothes, error) {

	var err error
	product_Clothes := []Product_Clothes{}
	err = db.Debug().Model(&Product_Clothes{}).Where("author_id = ?", uid).Limit(100).Order("created_at desc").Find(&product_Clothes).Error
	if err != nil {
		return &[]Product_Clothes{}, err
	}
	if len(product_Clothes) > 0 {
		for i, _ := range product_Clothes {
			err := db.Debug().Model(&Users{}).Where("id = ?", product_Clothes[i].AuthorID).Take(&product_Clothes[i].Author).Error
			if err != nil {
				return &[]Product_Clothes{}, err
			}
		}
	}
	return &product_Clothes, nil
}