package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alexeyco/simpletable"
	"github.com/razeenf/backup/logs"
)

func list(pwd string) {

	var log logs.BackupLog
	homeDir, _ := os.UserHomeDir()
	logDir := filepath.Join(homeDir, ".config/gdrive/log.json")
	log.Load(logDir)

	table := simpletable.New()
	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "Type"},
			{Align: simpletable.AlignCenter, Text: "Name"},
			{Align: simpletable.AlignCenter, Text: "Backed Up"},
			{Align: simpletable.AlignCenter, Text: "Last Backup"},
		},
	}
	var cells [][]*simpletable.Cell

	files, _ := os.ReadDir(pwd)
	for _, file := range files {

		var typee string
		if typee = "File"; file.IsDir() {
			typee = "Folder"
		}

		abs, _ := filepath.Abs(file.Name())
		fileInfo, _ := os.Stat(abs)

		backedUp := fmt.Sprint(string("\033[31m"), "NO", string("\033[0m"))
		backupTime := ""

		for _, v := range log {
			if strings.Contains(abs, v.DirName) {
				if fileInfo.ModTime().Before(v.BackupTime) {
					backedUp = fmt.Sprint(string("\033[32m"), "YES", string("\033[0m"))
					backupTime = v.BackupTime.Format("2006-01-02 03:04:05")
				} else if fileInfo.ModTime().After(v.BackupTime) {
					backedUp = fmt.Sprint(string("\033[33m"), "NEW", string("\033[0m"))
					backupTime = v.BackupTime.Format("2006-01-02 03:04:05")
				}
				break
			}
		}

		cells = append(cells, []*simpletable.Cell{
			{Text: typee},
			{Text: file.Name()},
			{Align: simpletable.AlignCenter, Text: backedUp},
			{Align: simpletable.AlignCenter, Text: backupTime},
		})
	}

	table.Body = &simpletable.Body{Cells: cells}
	table.Footer = &simpletable.Footer{Cells: []*simpletable.Cell{
		{Align: simpletable.AlignCenter, Span: 4, Text: "Directory: " + pwd},
	}}

	table.SetStyle(simpletable.StyleUnicode)
	table.Println()
}
