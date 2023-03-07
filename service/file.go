package service

import (
	"fmt"
	"os"
	"os/exec"
)

func DirCreateIfNotExist(dir string) {
	_, err := os.Stat(dir)
	if err != nil {
		os.MkdirAll(dir, os.ModePerm)
	}
}

func WriteFileAndCompress(targetDir string, fileName string, content []byte) {
	DirCreateIfNotExist(targetDir)
	// 將檔案寫入
	f, _ := os.OpenFile(targetDir+"/"+fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	defer f.Close()
	f.Write(content)
	returnData, err := exec.Command("convert", targetDir+"/"+fileName, "-quality", "80", targetDir+"/"+fileName).CombinedOutput()
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	fmt.Printf("returnData: %v\n", string(returnData))
}
