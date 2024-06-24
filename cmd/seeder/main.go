package main

import (
	"crypto/sha256"
	"encoding/csv"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/ochom/gutils/logs"
	"github.com/ochom/gutils/sql"
)

var threshold = 50_000_000

type Number struct {
	Id    int    `gorm:"primaryKey"`
	Hash  string `gorm:"unique"`
	Phone string
}

func init() {
	// init sqlite database
	cfg := sql.Config{
		Driver: sql.Sqlite,
		Url:    "data/seeder.db",
	}

	if err := sql.New(&cfg); err != nil {
		panic(err)
	}

	if err := sql.Conn().AutoMigrate(&Number{}); err != nil {
		panic(err)
	}
}

func main() {
	count := sql.Count[Number]()

	// clear database
	if err := clearData(count); err != nil {
		panic(err)
	}

	// seed data
	start := time.Now()
	if err := seed(count); err != nil {
		panic(err)
	}

	logs.Info("seeding took %s", time.Since(start))
}

func clearData(count int) error {
	start := time.Now()
	if count > threshold {
		return nil
	}

	if err := sql.Conn().Exec("DELETE FROM numbers").Error; err != nil {
		return err
	}

	logs.Info("clearing took %s", time.Since(start))
	return nil
}

func seed(count int) error {
	if count < threshold {
		// populateNumbers ...
		populateNumbers()
	}

	return writeNumbersToDatabase()
}

func populateNumbers() {
	// seed number for 70X
	for i := 700; i <= 709; i++ {
		createNumber(fmt.Sprintf("254%d", i), 1000000, "%06d")
	}

	// seed number for 71X
	for i := 710; i <= 719; i++ {
		createNumber(fmt.Sprintf("254%d", i), 1000000, "%06d")
	}

	// seed number for 72X
	for i := 720; i <= 729; i++ {
		createNumber(fmt.Sprintf("254%d", i), 1000000, "%06d")
	}

	// seed number between 740 and 743
	for i := 740; i <= 743; i++ {
		createNumber(fmt.Sprintf("254%d", i), 1000000, "%06d")
	}

	// seed more...
	prefixes := []int{745, 746, 748, 757, 758, 759}
	for _, p := range prefixes {
		createNumber(fmt.Sprintf("254%d", p), 1000000, "%06d")
	}

	// seed more...
	prefixes = []int{768, 769}
	for _, p := range prefixes {
		createNumber(fmt.Sprintf("254%d", p), 1000000, "%06d")
	}

	// seed 790...799
	for i := 790; i <= 799; i++ {
		createNumber(fmt.Sprintf("254%d", i), 1000000, "%06d")
	}

	// seed 11x
	for i := 110; i <= 119; i++ {
		createNumber(fmt.Sprintf("254%d", i), 1000000, "%06d")
	}
}

func createNumber(prefix string, max int, format string) {
	logs.Info("Seeding %s ...", prefix)
	start := time.Now()

	csvFile, err := os.OpenFile("numbers.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
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
		return
	}

	logs.Info("Seeding [%s] took %s", prefix, time.Since(start))
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

	jobs := make(chan string, threshold)
	results := make(chan Number, threshold)

	numWorkers := runtime.NumCPU()
	var wg sync.WaitGroup

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for phoneNumber := range jobs {
				results <- Number{
					Hash:  hashNumber(phoneNumber),
					Phone: phoneNumber,
				}
			}
		}()
	}

	for {
		row, err := reader.Read()
		if err != nil {
			break
		}

		jobs <- row[0]
	}

	close(jobs)

	wg.Wait()
	close(results)

	numbers := []Number{}
	for r := range results {
		numbers = append(numbers, r)
	}

	logs.Info("writing %d numbers to database", len(numbers))
	if err := sql.Conn().CreateInBatches(&numbers, 1000).Error; err != nil {
		return err
	}

	return nil
}
