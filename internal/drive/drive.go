package drive

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/omaressameldin/lunch-roulette/internal/env"
	drive "google.golang.org/api/drive/v3"
)

func DBFileName() string {
	dbName := env.GetDBName()

	return fmt.Sprintf("%s/%s.db", dbDirectory, dbName)
}

func UpdateDB() {
	if !canUseDrive() {
		return
	}

	fileName := DBFileName()
	srv := serviceAccount()
	var file *os.File
	file, err := os.Open(fileName)
	defer file.Close()
	if err != nil {
		return
	}

	// delete file if exists since update is not working
	defer deleteFile(srv, getFile(srv, fileName))

	f := &drive.File{Name: fileName, Capabilities: &drive.FileCapabilities{
		CanDownload: true,
		CanDelete:   true,
	}}
	_, err = srv.Files.Create(f).Media(file).Do()
	if err != nil {
		log.Fatalln(err)
	}
}

func GetDBFile() {
	if !canUseDrive() {
		return
	}

	fileName := DBFileName()
	srv := serviceAccount()

	f := getFile(srv, fileName)
	if f == nil {
		return
	}

	content, err := srv.Files.Get(f.Id).Download()
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(content.Body)
	if err != nil {
		log.Fatalf("An error occurred: %v\n", err)
	}
	ioutil.WriteFile(fileName, body, 0644)
}
