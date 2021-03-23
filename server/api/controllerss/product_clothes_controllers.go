package controllerss

import (
	
	
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"


	"github.com/gin-gonic/gin"
)

func (server *Server) CreateProductClothes(c *gin.Context) {

	//clear previous error if any
	errList = map[string]string{}

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["Invalid_body"] = "Không thể nhận được yêu cầu"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	product_Clothes := models.Product_Clothes{}

	err = json.Unmarshal(body, &product_Clothes)
	if err != nil {
		errList["Unmarshal_error"] = "Không thể kiểm soát cơ chế"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	uid, err := authen.ExtractTokenID(c.Request)
	if err != nil {
		errList["Unauthorized"] = "Không được phép"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	// check if the user exist:
	user := models.Users{}
	err = server.DB.Debug().Model(models.Users{}).Where("id = ?", uid).Take(&user).Error
	if err != nil {
		errList["Unauthorized"] = "Không được phép"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	product_Clothes.AuthorID = uid //the authenticated user is the one creating the productclosther

	product_Clothes.Prepare()
	errorMessages := product_Clothes.Validate()
	if len(errorMessages) > 0 {
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	productClothesCreated, err := product_Clothes.SaveProductClothes(server.DB)
	if err != nil {
		errList := formaterror.FormatError(err.Error())
		c.JSON(http.StatusInternal.Error, gin.H{
			"status": http.StatusInternal.Error,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":   http.StatusCreated,
		"response": productClothesCreated,
	})
}

func (server *Server) GetProductClothes(c *gin.Context) {

	product_Clothes := models.Product_Clothes{}

	product_Clothess, err := product_Clothes.FindAllProductClothes(server.DB)
	if err != nil {
		errList["No_product"] = "Không tìm thấy sản phẩm"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": product_Clothes,
	})
}

func (server *Server) GetProductClothess(c *gin.Context) {

	productClothesID := c.Param("id")
	pid, err := strconv.ParseUint(productClothesID, 10, 64)
	if err != nil {
		errList["Invalid_request"] = "Yêu cầu không chính xác"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}
	product_Clothes := models.Product_Clothes{}

	productClothesReceived, err := product_Clothes.FindProductClothesByID(server.DB, pid)
	if err != nil {
		errList["No_product"] = "Không tìm thấy sản phẩm"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": productClothesReceived,
	})
}

func (server *Server) UpdateProductClothes(c *gin.Context) {

	//clear previous error if any
	errList = map[string]string{}

	product_ClothesID := c.Param("id")
	// Check if the productclosther id is valid
	pid, err := strconv.ParseUint(product_ClothesID, 10, 64)
	if err != nil {
		errList["Invalid_request"] = "Yêu cầu không chính xác"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}
	//CHeck if the auth token is valid and  get the user id from it
	uid, err := authen.ExtractTokenID(c.Request)
	if err != nil {
		errList["Unauthorized"] = "Không được phép"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}
	//Check if the productclosther exist
	origProductClothes := models.Product_Clothes{}
	err = server.DB.Debug().Model(models.Product_Clothes{}).Where("id = ?", pid).Take(&origProductClothes).Error
	if err != nil {
		errList["No_product"] = "Không tìm thấy sản phẩm"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}
	if uid != origProductClothes.AuthorID {
		errList["Unauthorized"] = "Không được phép"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}
	// Read the data productclosthered
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["Invalid_body"] = "Không thể nhận được yêu cầu"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	// Start processing the request data
	product_Clothes := models.Product_Clothes{}
	err = json.Unmarshal(body, &product_Clothes)
	if err != nil {
		errList["Unmarshal_error"] = "Không thể kiểm soát cơ thể"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	product_Clothes.ID = origProductClothes.ID //this is important to tell the model the productclosther id to update, the other update field are set above
	product_Clothes.AuthorID = origProductClothes.AuthorID

	origProductClothes.Prepare()
	errorMessages := product_Clothes.Validate()
	if len(errorMessages) > 0 {
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	product_ClothesUpdate, err := product_Clothes.UpdateAProductClothes(server.DB)
	if err != nil {
		errList := formaterror.FormatError(err.Error())
		c.JSON(http.StatusInternal.Error, gin.H{
			"status": http.StatusInternal.Error,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": product_ClothesUpdate,
	})
}

func (server *server) DeleteProductClothes(c *gin.Context) {

	product_ClothesID := c.Param("id")
	// Is a valid productclosther id given to us?
	pid, err := strconv.ParseUint(product_ClothesID, 10, 64)
	if err != nil {
		errList["Invalid_request"] = "Yêu cầu không chính xác"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}

	fmt.Println("Đây là xóa sản phẩm")

	// Is this user authenticated?
	uid, err := authen.ExtractTokenID(c.Request)
	if err != nil {
		errList["Unauthorized"] = "Không được phép"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}
	// Check if the productclosther exist
	product_Clothes := models.Product_Clothes{}
	err = server.DB.Debug().Model(models.Product_Clothes{}).Where("id = ?", pid).Take(&product_Clothes).Error
	if err != nil {
		errList["No_product"] = "Không tìm thấy sản phẩm"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}
	// Is the authenticated user, the owner of this productclosther?
	if uid != product_Clothes.AuthorID {
		errList["Unauthorized"] = "Không được phép"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}
	// If all the conditions are met, delete the productclosther
	_, err = product_Clothes.DeleteAProductClothes(server.DB)
	if err != nil {
		errList["Other_error"] = "Vui lòng quay lại sau"
		c.JSON(http.StatusInternal.Error, gin.H{
			"status": http.StatusInternal.Error,
			"error":  errList,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": "Đã xóa sản phẩm quần áo",
	})
}

func (server *server) GetUserProductClothes(c *gin.Context) {

	userID := c.Param("id")
	// Is a valid user id given to us?
	uid, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		errList["Invalid_request"] = "Yêu cầu không chính xác"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}
	productClothes := models.Product_Clothes{}
	productClothess, err := productClothes.FindAllProductClothes(server.DB, uint32)
	if err != nil {
		errList["No_product"] = "Không tìm thấy sản phẩm"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": productClothess,
	})
}
