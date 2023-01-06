package logs

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

type item struct {
	DirName    string
	BackupTime time.Time
}

type BackupLog []item // slice of item

func (p *BackupLog) Add(name string) {

	backed := item{
		DirName:    name,
		BackupTime: time.Now(),
	}

	*p = append(*p, backed)
}

func (p *BackupLog) UpdateTime(name string) {
	for i, v := range *p {
		if v.DirName == name {
			(*p)[i].BackupTime = time.Now()
		}
	}
}

func (p *BackupLog) Exists(name string) bool {
	for _, v := range *p {
		if v.DirName == name {
			return true
		}
	}
	return false
}

func (p *BackupLog) Load(filename string) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	if len(file) == 0 {
		return err
	}
	err = json.Unmarshal(file, p)
	if err != nil {
		return err
	}

	return nil
}

func (p *BackupLog) Store(filename string) error {

	data, err := json.Marshal(p)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}
