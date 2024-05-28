package config

import (
	"log"
	"os"
)

func MakeFolder(folderName string) {
	path := "./" + folderName
	_, checkFolder := os.Stat(path)

	if checkFolder != nil {
		for checkFolder != nil {
			err := os.Mkdir(path, 0755)
			if err != nil {
				log.Println("Error creating directory")
			}
			_, checkFolder = os.Stat(path)
		}
	}
}
