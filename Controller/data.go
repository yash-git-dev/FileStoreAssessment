package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var db *sql.DB

type User struct {
	Name      string `json:"name"`
	Age       int    `json:"age"`
	Email     string `json:"email"`
	TimeStamp time.Time
}

func StoreData(c *gin.Context) {
	var user User

	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "error": "Failed to parse data"})
		return
	}

	user.TimeStamp = time.Now()

	userData, err := json.Marshal(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "error": "Failed to parse data"})
		return
	}

	filename := strings.ReplaceAll(user.Name, " ", "") + fmt.Sprintf("%d", time.Now().Unix()) + ".json"

	tmpDir := "tmp/astra"

	if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
		err := os.MkdirAll(tmpDir, 0755)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "error": "Internal server error"})
			return
		}
	}

	filePath := filepath.Join(tmpDir, filename)
	if err := ioutil.WriteFile(filePath, userData, 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "error": "Internal server error"})
		return
	}

	db, err = ConnectDB()
	if err != nil {
		log.Println(err.Error())
	}

	go createUser(filePath, db)

	c.JSON(http.StatusOK, gin.H{"message": "Data stored successfully"})
}

func createUser(filePath string, db *sql.DB) {
	var (
		user     User
		fileLock sync.RWMutex
	)

	fileLock.RLock()
	defer fileLock.RUnlock()

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("error reading from file: %v", err)
	}

	err = json.Unmarshal(data, &user)
	if err != nil {
		log.Printf("error unmarshalling data: %v", err)
	}

	query := "INSERT INTO `user_data` (name, age, email, created_at) VALUES (?, ?, ?, ?)"

	result, err := db.Exec(query, user.Name, user.Age, user.Email, user.TimeStamp)
	if err != nil {
		log.Printf("error executing query: %v", err)
	}

	rowsAff, err := result.RowsAffected()
	if rowsAff != 1 || err != nil {
		log.Println("no rows affected")
	} else {
		if err := os.Remove(filePath); err != nil {
			log.Printf("error deleting file: %v", err)
		}
	}

	defer db.Close()
}
