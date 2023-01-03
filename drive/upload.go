package drive

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/razeenf/backup/logs"
	"github.com/schollz/progressbar/v3"
	"google.golang.org/api/drive/v3"
)

func Backup(currentDir string) {
	//Step 1: Check if backup folder exists in drive and create if not exists
	srv, _ := getDriveService()
	exists, _ := srv.Files.List().Q("mimeType = 'application/vnd.google-apps.folder' and name = 'Backup Folder'").Do()
	var backupFolderId string
	if len(exists.Files) == 0 {
		folder, _ := createFolder(srv, "Backup Folder", "root")
		backupFolderId = folder.Id
	} else {
		backupFolderId = exists.Files[0].Id
	}

	//Step 2: Create subfolder in drive with name of current directory, if subfolder already exists, delete it and create new one
	part := strings.Split(currentDir, "\\")
	subfolderName := part[len(part)-2] + "/" + part[len(part)-1]
	exists, _ = srv.Files.List().Q("mimeType = 'application/vnd.google-apps.folder' and name = '" + subfolderName + "'").Do()

	if len(exists.Files) > 0 {
		for _, file := range exists.Files {
			srv.Files.Delete(file.Id).Do()
		}
	}
	dirFolder, _ := createFolder(srv, subfolderName, backupFolderId)

	//Step 3: Upload all files in current directory to drive
	files, err := os.ReadDir(currentDir)
	if err != nil {
		log.Fatal(err)
	}

	bar := progressbar.NewOptions((len(files)),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetWidth(60),
		progressbar.OptionSetDescription("Uploading Files:"),
		progressbar.OptionShowElapsedTimeOnFinish(),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[cyan]=[reset]",
			SaucerHead:    "[cyan]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))

	for _, file := range files {
		dir := filepath.Join(currentDir, file.Name())
		f, err := os.Open(dir)
		if err != nil {
			panic(fmt.Sprintf("cannot open file: %v", err))
		}
		defer f.Close()

		if file.IsDir() {
			subfolderUpload(srv, dirFolder.Id, dir)
		} else {
			createFile(srv, file.Name(), "application/octet-stream", f, dirFolder.Id)
		}
		bar.Add(1)
	}

	//Step 4: Update log file
	var log logs.BackupLog
	homeDir, _ := os.UserHomeDir()
	logDir := filepath.Join(homeDir, ".config/gdrive/log.json")

	if err := log.Load(logDir); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	subfolderName = part[len(part)-2] + "\\" + part[len(part)-1]

	if log.Exists(subfolderName) {
		log.UpdateTime(subfolderName)
	} else {
		log.Add(subfolderName)
	}
	log.Store(logDir)

	fmt.Println(string("\033[32m"), "\nBackup Complete at:", time.Now().Format("15:04:05 01-02-2006"), string("\033[0m"))
}

func subfolderUpload(srv *drive.Service, folderId string, subfolder string) {
	part := strings.Split(subfolder, "\\")
	folder, _ := createFolder(srv, part[len(part)-1], folderId)

	files, err := os.ReadDir(subfolder)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		dir := filepath.Join(subfolder, file.Name())
		if file.IsDir() {
			go subfolderUpload(srv, folder.Id, dir)
			continue
		}

		f, err := os.Open(dir)
		if err != nil {
			fmt.Printf("cannot open file: %v\n", file.Name())
			continue
		}
		defer f.Close()
		createFile(srv, file.Name(), "application/octet-stream", f, folder.Id)
	}

}

func createFolder(service *drive.Service, name string, parentId string) (*drive.File, error) {
	d := &drive.File{
		Name:     name,
		MimeType: "application/vnd.google-apps.folder",
		Parents:  []string{parentId},
	}

	file, err := service.Files.Create(d).Do()

	if err != nil {
		log.Println("Could not create dir: " + err.Error())
		return nil, err
	}

	return file, nil
}

func createFile(service *drive.Service, name string, mimeType string, content io.Reader, parentId string) (*drive.File, error) {
	f := &drive.File{
		MimeType: mimeType,
		Name:     name,
		Parents:  []string{parentId},
	}
	file, err := service.Files.Create(f).Media(content).Do()
	if err != nil {
		log.Println("Could not create file: " + err.Error())
		return nil, err
	}

	return file, nil
}
