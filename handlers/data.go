package handlers

import (
	"golang-assignment/config"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetEmployees(c *gin.Context) {
	employeeKeys, err := config.RedisClient.Keys(config.Ctx, "employee:*").Result()
	if err != nil {
		log.Println("Redis error:", err)
	} else if len(employeeKeys) > 0 {
		var result []map[string]string
		for _, key := range employeeKeys {
			employee, _ := config.RedisClient.HGetAll(config.Ctx, key).Result()
			result = append(result, employee)
		}
		c.JSON(http.StatusOK, result)
		return
	}

	rows, err := config.DB.Query("SELECT id, first_name, last_name, company_name, address, city, county, postal, phone, email, web FROM employees")
	if err != nil {
		log.Println("Database query error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query error"})
		return
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var id int
		var firstName, lastName, companyName, address, city, county, postal, phone, email, web string
		err := rows.Scan(&id, &firstName, &lastName, &companyName, &address, &city, &county, &postal, &phone, &email, &web)
		if err != nil {
			log.Println("Failed to scan row:", err)
			continue
		}
		result = append(result, map[string]interface{}{
			"id": id, "first_name": firstName, "last_name": lastName, "company_name": companyName,
			"address": address, "city": city, "county": county, "postal": postal, "phone": phone,
			"email": email, "web": web,
		})
	}

	if len(result) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No records found"})
	} else {
		c.JSON(http.StatusOK, result)
	}
}

func UpdateEmployee(c *gin.Context) {
	id := c.Param("id")
	var employee map[string]string
	if err := c.ShouldBindJSON(&employee); err != nil {
		log.Println("Failed to bind JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Update MySQL
	query := `UPDATE employees SET first_name=?, last_name=?, company_name=?, address=?, city=?, county=?, postal=?, phone=?, email=?, web=? WHERE id=?`
	if _, err := config.DB.Exec(query, employee["first_name"], employee["last_name"], employee["company_name"], employee["address"], employee["city"], employee["county"], employee["postal"], employee["phone"], employee["email"], employee["web"], id); err != nil {
		log.Println("Failed to update record in MySQL:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update record"})
		return
	}

	// Update Redis
	key := "employee:" + id
	if err := config.RedisClient.HSet(config.Ctx, key, employee).Err(); err != nil {
		log.Println("Failed to update cache in Redis:", err)
	}
	config.RedisClient.Expire(config.Ctx, key, 5*time.Minute)

	c.JSON(http.StatusOK, gin.H{"message": "Record updated successfully"})
}

func DeleteEmployee(c *gin.Context) {
	id := c.Param("id")

	// Delete from MySQL
	query := `DELETE FROM employees WHERE id=?`
	if _, err := config.DB.Exec(query, id); err != nil {
		log.Println("Failed to delete record in MySQL:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete record"})
		return
	}

	// Delete from Redis
	key := "employee:" + id
	if err := config.RedisClient.Del(config.Ctx, key).Err(); err != nil {
		log.Println("Failed to delete cache in Redis:", err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Record deleted successfully"})
}
