package handlers

import (
	"golang-assignment/utils"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		log.Println("Failed to get form file:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get form file"})
		return
	}

	uploadPath := filepath.Join("uploads", file.Filename)
	if err := os.MkdirAll(filepath.Dir(uploadPath), os.ModePerm); err != nil {
		log.Println("Failed to create upload directory:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	if err := c.SaveUploadedFile(file, uploadPath); err != nil {
		log.Println("Failed to save uploaded file:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save uploaded file"})
		return
	}

	go utils.ProcessExcelFile(uploadPath)

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully"})
}
