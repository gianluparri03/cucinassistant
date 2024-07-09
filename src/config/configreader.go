package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

func Read(path string) {
	var data []byte
	var err error

	// Reads the file
	if data, err = os.ReadFile(path); err != nil {
		log.Fatal("ERR: " + err.Error())
	}

	// Parses it
	if err = yaml.Unmarshal(data, &Runtime); err != nil {
		log.Fatal("ERR: " + err.Error())
	}
}
