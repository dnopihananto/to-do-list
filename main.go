package main

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `gorm:"type:varchar(100)"`
	Position string `gorm:"type:varchar(100)"`
}

type ToDo struct {
	gorm.Model
	Title   string `gorm:"type:varchar(100)"`
	Desc    string `gorm:"type:text"`
	UserId  int
	OrderId int64
	Status  string `gorm:"type:text"`
}

var DB *gorm.DB

func main() {
	dsn := "root:@tcp(127.0.0.1:3306)/learn_gin_todo_list?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect")
	}
	DB.AutoMigrate(&User{}, &ToDo{})

	router := gin.Default()

	v1 := router.Group("/api/v1/")
	{
		users := v1.Group("/user")
		{
			users.GET("/", getUsers)
			users.GET("/:id", getUser)
			users.POST("/", postUser)
			users.POST("/:id", updateUser)
			users.DELETE("/:id", deleteUser)
		}

		todos := v1.Group("/to-do")
		{
			todos.GET("/", getToDoList)
			todos.GET("/:id", getToDo)
			todos.POST("/", postToDo)
			todos.POST("/:id", updateToDo)
			todos.DELETE("/:id", deleteToDo)
			todos.POST("/:id/change-order-id", changeOrderId)
		}
	}

	router.Run()
}

func getUsers(c *gin.Context) {

	items := []User{}
	DB.Find(&items)

	c.JSON(200, gin.H{
		"status":  "00",
		"message": "berhasil ambil data semua user",
		"data":    items,
	})
}

func getUser(c *gin.Context) {
	id := c.Param("id")

	var item User

	dbRresult := DB.First(&item, "id = ?", id)
	if errors.Is(dbRresult.Error, gorm.ErrRecordNotFound) {
		c.JSON(404, gin.H{
			"status":  "error",
			"message": "record not found",
			"id":      id,
		})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{
		"status":  "berhasil",
		"message": "berhasil ambil data user",
		"data":    item,
	})
}

func postUser(c *gin.Context) {
	item := User{
		Name:     c.PostForm("name"),
		Position: c.PostForm("position"),
	}

	DB.Create(&item)

	c.JSON(200, gin.H{
		"status": "berhasil simpan data",
		"data":   item,
	})
}

func updateUser(c *gin.Context) {
	id := c.Param("id")

	var item User

	dbRresult := DB.First(&item, "id = ?", id)
	if errors.Is(dbRresult.Error, gorm.ErrRecordNotFound) {
		c.JSON(404, gin.H{
			"status":  "error",
			"message": "record not found",
		})
		c.Abort()
		return
	}

	item.Name = c.PostForm("name")
	item.Position = c.PostForm("position")
	DB.Save(&item)

	c.JSON(200, gin.H{
		"status": "berhasil ubah data",
		"data":   item,
	})
}

func deleteUser(c *gin.Context) {
	id := c.Param("id")

	var item User

	dbRresult := DB.First(&item, "id = ?", id)
	if errors.Is(dbRresult.Error, gorm.ErrRecordNotFound) {
		c.JSON(404, gin.H{
			"status":  "error",
			"message": "record not found",
		})
		c.Abort()
		return
	}
	DB.Delete(&item)

	c.JSON(200, gin.H{
		"status": "berhasil hapus data",
		"data":   item,
	})
}

func getToDoList(c *gin.Context) {

	items := []ToDo{}
	DB.Find(&items)

	c.JSON(200, gin.H{
		"status":  "00",
		"message": "berhasil ambil data semua to do",
		"data":    items,
	})
}

func getToDo(c *gin.Context) {
	id := c.Param("id")

	var item ToDo

	dbRresult := DB.First(&item, "id = ?", id)
	if errors.Is(dbRresult.Error, gorm.ErrRecordNotFound) {
		c.JSON(404, gin.H{
			"status":  "error",
			"message": "record not found",
		})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{
		"status":  "berhasil",
		"message": "berhasil ambil data to do",
		"data":    item,
	})
}

func postToDo(c *gin.Context) {
	var totalTodoPerStatus int64
	DB.Model(&ToDo{}).Where("status = ?", c.PostForm("status")).Count(&totalTodoPerStatus)
	totalTodoPerStatus = totalTodoPerStatus + 1
	userId, _ := strconv.Atoi(c.PostForm("user_id"))
	item := ToDo{
		Title:   c.PostForm("title"),
		Desc:    c.PostForm("desc"),
		UserId:  userId,
		OrderId: totalTodoPerStatus,
		Status:  c.PostForm("status"),
	}

	DB.Create(&item)

	c.JSON(200, gin.H{
		"status": "berhasil simpan data",
		"data":   item,
	})
}

func updateToDo(c *gin.Context) {
	id := c.Param("id")

	var item ToDo

	dbRresult := DB.First(&item, "id = ?", id)
	if errors.Is(dbRresult.Error, gorm.ErrRecordNotFound) {
		c.JSON(404, gin.H{
			"status":  "error",
			"message": "record not found",
		})
		c.Abort()
		return
	}

	userId, _ := strconv.Atoi(c.PostForm("user_id"))
	item.Title = c.PostForm("title")
	item.Desc = c.PostForm("desc")
	item.UserId = userId
	item.Status = c.PostForm("status")
	DB.Save(&item)

	c.JSON(200, gin.H{
		"status": "berhasil ubah data",
		"data":   item,
	})
}

func deleteToDo(c *gin.Context) {
	id := c.Param("id")

	var item ToDo

	dbRresult := DB.First(&item, "id = ?", id)
	if errors.Is(dbRresult.Error, gorm.ErrRecordNotFound) {
		c.JSON(404, gin.H{
			"status":  "error",
			"message": "record not found",
		})
		c.Abort()
		return
	}
	DB.Delete(&item)

	c.JSON(200, gin.H{
		"status": "berhasil hapus data",
		"data":   item,
	})
}

func changeOrderId(c *gin.Context) {
	id := c.Param("id")

	var item ToDo

	dbRresult := DB.First(&item, "id = ?", id)
	if errors.Is(dbRresult.Error, gorm.ErrRecordNotFound) {
		c.JSON(404, gin.H{
			"status":  "error",
			"message": "record not found",
		})
		c.Abort()
		return
	}

	orderId, _ := strconv.ParseInt(c.PostForm("order_id"), 10, 64)
	item.OrderId = orderId
	DB.Save(&item)

	DB.Exec("UPDATE to_dos SET order_id = order_id + 1 WHERE status = ? AND id <> ? AND deleted_at IS NOT NULL", item.Status, item.ID)

	c.JSON(200, gin.H{
		"status": "berhasil ubah data",
		"data":   item,
	})
}
