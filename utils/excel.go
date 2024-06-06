package utils

import (
	"golang-assignment/config"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/xuri/excelize/v2"
)

func ProcessExcelFile(filePath string) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		log.Println("Failed to open Excel file:", err)
		return
	}

	sheetNames := f.GetSheetList()
	log.Println("Available sheets:", sheetNames)
	if len(sheetNames) == 0 {
		log.Println("No sheets found in the Excel file")
		return
	}

	sheetName := sheetNames[0]
	rows, err := f.GetRows(sheetName)
	if err != nil {
		log.Println("Failed to get rows:", err)
		return
	}

	headers := rows[0]
	log.Println("Headers found:", headers)

	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 10 {
			log.Println("Incomplete row:", row)
			continue
		}

		firstName := row[0]
		lastName := row[1]
		companyName := row[2]
		address := row[3]
		city := row[4]
		county := row[5]
		postal := row[6]
		phone := row[7]
		email := row[8]
		web := row[9]

		query := `INSERT INTO employees (first_name, last_name, company_name, address, city, county, postal, phone, email, web) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
		if _, err := config.DB.Exec(query, firstName, lastName, companyName, address, city, county, postal, phone, email, web); err != nil {
			log.Println("Failed to insert record into MySQL:", err)
			continue
		}

		employee := map[string]interface{}{
			"first_name":   firstName,
			"last_name":    lastName,
			"company_name": companyName,
			"address":      address,
			"city":         city,
			"county":       county,
			"postal":       postal,
			"phone":        phone,
			"email":        email,
			"web":          web,
		}
		key := "employee:" + strconv.Itoa(i)
		if err := config.RedisClient.HSet(config.Ctx, key, employee).Err(); err != nil {
			log.Println("Failed to cache record in Redis:", err)
		}
		config.RedisClient.Expire(config.Ctx, key, 5*time.Minute)
	}
	log.Println("Excel file processed successfully")

	if err := os.Remove(filePath); err != nil {
		log.Println("Failed to delete the Excel file:", err)
	} else {
		log.Println("Excel file deleted successfully")
	}
}
