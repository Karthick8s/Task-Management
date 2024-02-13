package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// var db *gorm.DB

type Priority string

const (
	Low    Priority = "low"
	Medium Priority = "medium"
	High   Priority = "high"
)

type errorMsg string

const (
	InvalidPayload    errorMsg = "Invalid payload"
	ErrorCreatingTask errorMsg = "Error Creating Task"
	ErrorGettingTasks errorMsg = "Error getting all tasks"
	InvalidTaskID     errorMsg = "Invalid Task ID"
	TaskNotFound      errorMsg = "Task Not found"
	UnableToUpdate    errorMsg = "Unable to update the task"
)

type Task struct {
	TaskID      uint64    `gorm:"primaryKey;autoIncrement" json:"task_id"`
	TaskName    string    `json:"task_name" gorm:"not null"`
	Assignee    uint64    `json:"assignee"`
	Priority    Priority  `json:"priority"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"startdate"`
	EndDate     time.Time `json:"enddate"`
}

type User struct {
	UserID   uint64 `gorm:"primaryKey;autoIncrement" json:"userid"`
	UserName string `json:"username"`
	Email    string `json:"email"`
	IsActive bool   `json:"isactive" gorm:"default:true"`
}

var db *gorm.DB

// DB connections
func initDB() {
	var err error

	dsn := "   "

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database ", err)
		return
	}

	// AutoMigrate will create the tasks table based on the Task struct
	err = db.AutoMigrate(&Task{})
	if err != nil {
		log.Fatal("Failed to migrate database", err)
		return
	}

	// err = db.AutoMigrate(&User{})
	// if err != nil {
	// 	log.Fatal("Failed to migrate User", err)
	// 	return
	// }

}

func main() {
	initDB()

	r := gin.Default()

	//tasks
	r.POST("/tasks", createTask)
	r.GET("/getTasks", getTask)
	r.GET("/getTask/:task_id", getTaskbyID)
	r.PUT("/task/:task_id", updateTask)
	r.DELETE("/task/:task_id", deleteTask)

	//users
	r.POST("/users/create", createUser)
	r.GET("/users", getAllUsers)
	r.GET("/getByID/:user_id", getUserByID)
	r.DELETE("/users/:user_id", deleteUser)
	r.PUT("/users/:user_id/update", updateUser)

	fmt.Println("Server  is running on http://localhost;8080")
	r.Run(":0808")
}

func createTask(c *gin.Context) {
	var newTask Task

	if err := c.ShouldBindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, InvalidPayload)
		return
	}

	// Create Task in the database
	if err := db.Debug().Create(&newTask).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorCreatingTask)
	}

	c.JSON(http.StatusCreated, newTask)
}

func getTask(c *gin.Context) {
	var tasks []Task

	query := "SELECT * FROM tasks"
	if err := db.Debug().Raw(query).Scan(&tasks).
		Error; err != nil {

		c.JSON(http.StatusInternalServerError, ErrorGettingTasks)
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func getTaskbyID(c *gin.Context) {
	var task Task

	taskID, err := strconv.Atoi(c.Param("task_id"))
	fmt.Println("taskId", taskID)
	if err != nil {
		c.JSON(http.StatusBadRequest, InvalidTaskID)
		return
	}

	if err := db.Debug().First(&task, uint64(taskID)).Error; err != nil {
		c.JSON(http.StatusNotFound, TaskNotFound)
		return
	}

	c.JSON(http.StatusOK, task)
}

func updateTask(c *gin.Context) {
	fmt.Println("Update Task")

	taskID, err := strconv.Atoi(c.Param("task_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, InvalidTaskID)
	}

	var updatedTask Task
	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, InvalidPayload)
		return
	}

	//Update in database
	if err := db.Debug().Model(&Task{}).
		Where("task_id", taskID).
		Updates(&updatedTask).Error; err != nil {
		c.JSON(http.StatusInternalServerError, UnableToUpdate)
		return
	}
	c.JSON(http.StatusOK, updatedTask)
}

func deleteTask(c *gin.Context) {
	fmt.Println("Delete Task")

	taskID, err := strconv.Atoi(c.Param("task_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, InvalidTaskID)
		return
	}

	if err := db.Debug().Delete(&Task{}, taskID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, "Failed deleted task")
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func createUser(c *gin.Context) {
	fmt.Println("Create User")

	var user User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, "Error while getting user payload")
		return
	}

	result := db.Debug().Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, "Failed Error while creating")
		return
	}

	c.JSON(http.StatusCreated, user)
}

func getAllUsers(c *gin.Context) {
	fmt.Println("Get users ")

	var users []User
	query := "SELECT * FROM users"
	if err := db.Raw(query).Debug().Scan(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, "Error getting users")
		return
	}

	c.JSON(http.StatusOK, users)
}

func getUserByID(c *gin.Context) {
	fmt.Println("---Get user by userID---")

	var user User
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid user ID")
		return
	}

	if err := db.First(&user, userID).Debug().Error; err != nil {
		c.JSON(http.StatusNotFound, "TaskID not found")
		return
	}

	c.JSON(http.StatusOK, user)
}

func deleteUser(c *gin.Context) {
	fmt.Println("---Delete User---")

	userID := c.Param("user_id")

	if err := db.First(&User{}, userID).Debug().Error; err != nil {
		c.JSON(http.StatusNotFound, "TaskID not found")
		return
	}

	if err := db.Debug().Delete(&User{}, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, "Failed to delete the user")
		return
	}

	c.JSON(http.StatusOK, "Deleted successfully")
}

func updateUser(c *gin.Context) {
	fmt.Println("---update user---")

	userID := c.Param("user_id")

	var user User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := db.Model(&User{}).Where("user_id", userID).Updates(&user)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, "Error while updating user")
		return
	}

	c.JSON(http.StatusAccepted, user)
}
