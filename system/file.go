package system

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

func WriteToFile(filechannel chan []string, fileName string) {
	datetime := strings.Replace(time.Now().Format("2006-01-02"), "-", "", -1)
	// format time to yyy_mm_dd

	fName := fmt.Sprintf("%v_%v.txt", fileName, datetime)

	file, err := os.Create(fmt.Sprintf("arquivos/%v", fName))
	if err != nil {
		log.Fatal("Error creating file: " + err.Error())
	}

	defer file.Close()

	for row := range filechannel {
		line := ""
		iLen := len(row)

		for i, col := range row {

			col = strings.TrimSpace(col)
			if col == "<nil>" {
				col = ""
			}
			if i+1 < iLen {
				col = col + "|"
			}

			line += col
		}

		_, err := file.WriteString(line + "\n")
		if err != nil {
			log.Fatal("Error writing line to file: " + err.Error())
		}
	}
}

func ZipFiles(files []string) error {
	fileZipName := fmt.Sprintf("./enviados/%v.zip", time.Now().Format("2006-01-02"))
	newZipFile, err := os.Create(fileZipName)
	if err != nil {
		return err
	}

	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	for _, file := range files {
		if err := addFileToZip(zipWriter, file); err != nil {
			return err
		}

	}
	return nil
}

func addFileToZip(zipWriter *zip.Writer, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	defer file.Close()

	info, err := file.Stat()

	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	fmt.Println(strings.Split(filename, "/"))
	header.Method = zip.Deflate
	header.Name = strings.Split(filename, "/")[2]

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, file)
	return err
}
