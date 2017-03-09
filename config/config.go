package config

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"
)

var (
	appConfigPath string
)

func getCurrentDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Panic(err)
	}
	return dir
}

func fileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
func init() {
	currentDir := getCurrentDir()
	appConfigPath = filepath.Join(currentDir, "config", "config.json")
}

func parseConfig(v interface{}) {
	if !fileExist(appConfigPath) {
		log.Panicf("config file [%s] not found", appConfigPath)
		return
	}
	f, _ := os.Open(appConfigPath)
	defer f.Close()
	reader := json.NewDecoder(f)
	for {
		if err := reader.Decode(v); err == io.EOF {
			break
		} else if err != nil {
			log.Panic(err)
		}
	}
}

func LoadConfig(v interface{}) {
	parseConfig(v)
}
