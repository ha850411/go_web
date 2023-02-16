package service

import (
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
	// 將檔案寫入
	f, _ := os.OpenFile(targetDir+"/"+fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	defer f.Close()
	f.Write(content)
	exec.Command("convert", "-quality", "80", targetDir+"/"+fileName, targetDir+"/"+fileName).Run()
}
