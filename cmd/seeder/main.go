package main

import (
	"crypto/sha256"
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/ochom/gutils/logs"
	"github.com/ochom/gutils/sql"
)

type Number struct {
	Id    int    `gorm:"primaryKey"`
	Hash  string `gorm:"unique"`
	Phone string
}

func init() {
	// init sqlite database
	cfg := sql.Config{
		DatabaseType: sql.Sqlite,
		Url:          "seeder.db",
	}

	if err := sql.New(&cfg); err != nil {
		panic(err)
	}

	if err := sql.Conn().AutoMigrate(&Number{}); err != nil {
		panic(err)
	}
}

func main() {
	// clear database
	if err := sql.Conn().Exec("DELETE FROM numbers").Error; err != nil {
		panic(err)
	}

	// seed data
	start := time.Now()
	if err := seed(); err != nil {
		panic(err)
	}

	logs.Info("seeding took %s", time.Since(start))
}

func seed() error {
	// populateNumbers ...
	populateNumbers()

	// create possible number of the format with format 25411xxxxxxx
	createNumber("25411", 10000000, "%07d")
	return writeNumbersToDatabase()
}

func populateNumbers() {
	// seed number between 701 and 729
	for i := 701; i <= 729; i++ {
		createNumber(fmt.Sprintf("254%d", i), 1000000, "%06d")
	}

	// seed number between 740 and 743
	for i := 740; i <= 743; i++ {
		createNumber(fmt.Sprintf("254%d", i), 1000000, "%06d")
	}

	// seed more...
	prefixes := []int{45, 46, 48, 57, 58, 59}
	for _, p := range prefixes {
		createNumber(fmt.Sprintf("2547%d", p), 1000000, "%06d")
	}

	// seed more...
	prefixes = []int{68, 69}
	for _, p := range prefixes {
		createNumber(fmt.Sprintf("2547%d", p), 1000000, "%06d")
	}

	// seed 90...99
	for i := 90; i <= 99; i++ {
		createNumber(fmt.Sprintf("2547%d", i), 1000000, "%06d")
	}
}

func createNumber(prefix string, max int, format string) error {
	logs.Info("Seeding %s ...", prefix)
	start := time.Now()

	csvFile, err := os.Create("numbers.csv")
	if err != nil {
		return err
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	records := [][]string{}
	for i := range max {
		number := fmt.Sprintf("%s%s", prefix, fmt.Sprintf(format, i))
		records = append(records, []string{number})
	}

	if err := writer.WriteAll(records); err != nil {
		return err
	}

	logs.Info("Seeding [%s] took %s", prefix, time.Since(start))
	return nil
}

func hashNumber(num string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(num)))
}

func writeNumbersToDatabase() error {
	file, err := os.Open("numbers.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	chunkSize := 1000
	chunk := make([][]string, 0, chunkSize)
	for {
		// Read up to chunkSize rows at a time
		for i := 0; i < chunkSize; i++ {
			row, err := reader.Read()
			if err != nil {
				break // End of file
			}
			chunk = append(chunk, row)
		}

		if err := saveChunk(chunk); err != nil {
			break
		}

		chunk = chunk[:0]
	}

	return nil
}

// func hashChunk(chunk [][]string) {
// 	numbers := []Number{}
// 	for _, row := range chunk {
// 		numbers = append(numbers, Number{
// 			Hash:  hashNumber(row[0]),
// 			Phone: row[0],
// 		})
// 	}

// 	if len(numbers) == 0 {
// 		return
// 	}

// 	csvFile, err := os.OpenFile("numbers_hashed.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer csvFile.Close()

// 	writer := csv.NewWriter(csvFile)
// 	defer writer.Flush()

// 	records := [][]string{}
// 	for _, num := range numbers {
// 		records = append(records, []string{num.Phone, num.Hash})
// 	}

// 	if err := writer.WriteAll(records); err != nil {
// 		panic(err)
// 	}

// 	logs.Info("Hashed %d numbers ...", len(numbers))
// }

func saveChunk(chunk [][]string) error {
	numbers := []Number{}
	for _, row := range chunk {
		numbers = append(numbers, Number{
			Hash:  hashNumber(row[0]),
			Phone: row[0],
		})
	}

	if err := sql.Conn().CreateInBatches(&numbers, 1000).Error; err != nil {
		return err
	}

	logs.Info("Saved %d numbers ...", len(numbers))
	return nil
}
