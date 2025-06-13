package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type Folder struct {
	FolderName     string   `json:"FolderName"`
	FileExtensions []string `json:"FileExtensions"`
}

type Folders struct {
	Folders []Folder `json:"Folders"`
}

func main() {

	dirPath := "/source"
	targetPathRoot := "/dest"

	files, err := os.ReadDir(dirPath)
	if err != nil {
		panic(err)
	}
	targetPathRootCreated := createDirectory(targetPathRoot)
	if targetPathRootCreated == false {
		fmt.Printf("Could not create directory %s", targetPathRoot)
		panic("Error")
	}

	folders, err := getFolderSettingsByJson("./folders.json")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		fileName := file.Name()
		fullPath := filepath.Join(dirPath, fileName)
		if file.IsDir() {
			continue
		}
		if fileName == ".DS_STORE" || fileName[0] == '.' {
			continue
		}

		moveFile(fileName, fullPath, folders, targetPathRoot)
	}

}

func getFolderSettingsByJson(pathToJson string) ([]Folder, error) {
	var wrapper Folders

	jsonFile, err := os.Open(pathToJson)
	if err != nil {
		return nil, fmt.Errorf("cannot open JSON file: %w", err)
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, fmt.Errorf("error reading JSON file: %w", err)
	}

	if err := json.Unmarshal(byteValue, &wrapper); err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}
	return wrapper.Folders, nil
}

func moveFile(rawFileName string, source string, mapFolder []Folder, dest string) {
	fileSplit := strings.Split(rawFileName, ".")
	fileExt := fileSplit[len(fileSplit)-1]
	for _, value := range mapFolder {
		if slices.Contains(value.FileExtensions, fileExt) {
			destination := dest + "/" + value.FolderName
			fmt.Println(destination)
			newDirectoryCreated := createDirectory(destination)
			if newDirectoryCreated == false {
				log.Fatal("Failed to create directory: $s", destination)
				break
			}
			destPath := filepath.Join(destination, rawFileName)
			err := os.Rename(source, destPath)
			if err != nil {
				log.Fatal(err)
				break
			}
			fmt.Printf("File %s has been created sucessfully \n", destPath)
		}
	}

}

func createDirectory(directory string) bool {
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		if err := os.MkdirAll(directory, os.ModePerm); err != nil {
			return false
		}
	}
	return true
}
