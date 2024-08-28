package main

import (
	"database/sql"
	"export_360/system"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/microsoft/go-mssqldb"
)

const (
	dir = "./arquivos"
)

func main() {
	start := time.Now()

	db, config := initialize()
	defer db.Close()

	processViews(db, config.Views)
	uploadFilesToS3(dir, config)
	compressFiles(dir)
	clearFolder(dir)

	elapsed := time.Since(start)
	fmt.Printf("Execution time: %s\n", elapsed)
}

func initialize() (*sql.DB, *system.Config) {
	db, err := system.ConnectToDb()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	config, err := system.ReadConfigFile()
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	return db, &config
}

func processViews(db *sql.DB, views []string) {
	for _, view := range views {
		log.Printf("Processing view: %s", view)
		dbchannel := make(chan []string, 500)
		go system.WriteToFile(dbchannel, view)
		system.GetDbData(db, dbchannel, view)
		close(dbchannel)
	}
}

func uploadFilesToS3(directory string, config *system.Config) {
	files, err := os.ReadDir(directory)
	if err != nil {
		log.Fatalf("Error reading directory: %v", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filepath := fmt.Sprintf("%s/%s", directory, file.Name())
		if err := system.Upload(filepath, config); err != nil {
			log.Fatalf("Error uploading file %s: %v", file.Name(), err)
		}
		log.Printf("File %s uploaded successfully", file.Name())
	}
}

func compressFiles(directory string) {
	files, err := os.ReadDir(directory)
	if err != nil {
		log.Fatalf("Error reading directory: %v", err)
	}

	var fileList []string
	for _, file := range files {
		fileList = append(fileList, fmt.Sprintf("%s/%s", directory, file.Name()))
	}

	if err := system.ZipFiles(fileList); err != nil {
		log.Fatalf("Error compressing files: %v", err)
	}
	fmt.Println("Files compressed successfully")
}

func clearFolder(directory string) {
	files, err := os.ReadDir(directory)
	if err != nil {
		log.Fatalf("Error reading directory: %v", err)
	}

	for _, file := range files {
		if err := os.Remove(fmt.Sprintf("%s/%s", directory, file.Name())); err != nil {
			log.Fatalf("Error removing file %s: %v", file.Name(), err)
		}
	}
}
