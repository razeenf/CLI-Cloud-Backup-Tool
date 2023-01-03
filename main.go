package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/razeenf/backup/drive"
)

func main() {
	listPtr := flag.Bool("list", false, "Lists all files in the current directory and whether they have been backed up or not.")
	backupPtr := flag.Bool("this", false, "Backs up the current directory.")

	flag.Parse()
	currentDir, _ := os.Getwd()

	if *listPtr {
		list(currentDir)
	}

	if *backupPtr {
		fmt.Println("Preparing to Backup Current Directory...")
		drive.Backup(currentDir)
	}
}
