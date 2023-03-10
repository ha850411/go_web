package service

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
)

func DirCreateIfNotExist(dir string) {
	_, err := os.Stat(dir)
	if err != nil {
		os.MkdirAll(dir, os.ModePerm)
	}
}

func WriteFile(targetDir string, fileName string, content []byte) {

	DirCreateIfNotExist(targetDir)
	// 將檔案寫入
	f, _ := os.OpenFile(targetDir+"/"+fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	defer f.Close()
	f.Write(content)
}

func WriteFileAndCompress(targetDir string, fileName string, content []byte) {

	DirCreateIfNotExist(targetDir)
	// 將檔案寫入
	f, _ := os.OpenFile(targetDir+"/"+fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	defer f.Close()
	f.Write(content)
	fi, err := os.Stat(targetDir + "/" + fileName)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	size := fi.Size()
	compressSize := 1024 * 1024 * 5 // 5MB 以上要壓縮
	if int(size) > compressSize {
		fmt.Println("檔案 > 5MB, 執行壓縮")
		contentType := http.DetectContentType(content)
		fmt.Printf("contentType: %v\n", contentType)
		if contentType == "image/jpeg" {
			returnData, err := exec.Command("convert", "-sample", "50%x50%", targetDir+"/"+fileName, targetDir+"/"+fileName).CombinedOutput()
			if err != nil {
				fmt.Printf("err: %v\n", err)
			}
			fmt.Printf("returnData: %v\n", string(returnData))
		} else {
			returnData, err := exec.Command("convert", targetDir+"/"+fileName, "-quality", "80", targetDir+"/"+fileName).CombinedOutput()
			if err != nil {
				fmt.Printf("err: %v\n", err)
			}
			fmt.Printf("returnData: %v\n", string(returnData))
		}
	}
}
