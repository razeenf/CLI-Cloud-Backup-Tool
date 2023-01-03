package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alexeyco/simpletable"
	"github.com/razeenf/backup/logs"
)

func list(pwd string) {

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

		var log logs.BackupLog
		homeDir, _ := os.UserHomeDir()
		logDir := filepath.Join(homeDir, ".config/gdrive/log.json")
		log.Load(logDir)

		abs, _ := filepath.Abs(file.Name())
		fileInfo, _ := os.Stat(abs)

		backedUp := fmt.Sprint(string("\033[31m"), "NO", string("\033[0m"))
		backupTime := ""

		for _, v := range log {
			timeString := v.BackupTime
			BUT, _ := time.Parse("2006-01-02 03:04:05", timeString)

			if strings.Contains(abs, v.DirName) && fileInfo.ModTime().Before(BUT) {
				backedUp = fmt.Sprint(string("\033[32m"), "YES", string("\033[0m"))
				backupTime = timeString
			} else if strings.Contains(abs, v.DirName) && fileInfo.ModTime().After(BUT) {
				backedUp = fmt.Sprint(string("\033[33m"), "NEW", string("\033[0m"))
				backupTime = "Modified Since"
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
