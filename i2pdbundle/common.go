package i2pd

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type fsi interface {
	IsDir() bool
	Readdir(int) ([]os.FileInfo, error)
	Open(name string) (http.File, error)
}

func FindAllDirectories(filesystem fsi) ([]string, error) {
	if filesystem.IsDir() {
		filelist, err := filesystem.Readdir(0)
		if err != nil {
			return nil, err
		}
		var rlist []string
		for index, file := range filelist {
			add := true
			for _, dir := range rlist {
				if dir == filepath.Dir(file.Name()) {
					add = false
				}
			}
			if add {
				rlist = append(rlist, filepath.Dir(file.Name()))
				log.Println(index, filepath.Dir(file.Name()))
			}
		}
	}
	return nil, nil
}

func FindAllFiles(filesystem fsi) ([]string, error) {
	if filesystem.IsDir() {
		filelist, err := filesystem.Readdir(0)
		if err != nil {
			return nil, err
		}
		var rlist []string
		for index, file := range filelist {
			if !file.IsDir() {
				rlist = append(rlist, file.Name())
				log.Println(index, file.Name())
			}
		}
	}
	return nil, nil
}

func WriteAllFiles(filesystem fsi, unpackdir string) error {
	dirs, err := FindAllDirectories(filesystem)
	if err != nil {
		return fmt.Errorf("Directory Discovery Error, %s", err)
	}
	for _, dir := range dirs {
		os.MkdirAll(unpackdir+dir, 0755)
		if err != nil {
			return fmt.Errorf("Directory Discovery Error, %s", err)
		}
	}
	if filesystem.IsDir() {
		log.Println("Found a directory, preparing to start loop")
		if filelist, err := filesystem.Readdir(0); err == nil {
			log.Println("Starting loop")
			for index, fi := range filelist {
				if file, err := filesystem.Open(fi.Name()); err == nil {
					if !fi.IsDir() {
						var buf []byte
						if _, err := file.Read(buf); err == nil {
							log.Println(index, fi.Name())
							if err := ioutil.WriteFile(unpackdir+"/"+fi.Name(), buf, fi.Mode()); err != nil {
								return fmt.Errorf("Write file error", err)
							}
							log.Println("Wrote file", fi.Name())
						} else {
							return fmt.Errorf("Read Error: %s, %s", fi.Name(), err)
						}
					}
					file.Close()
				} else {
					return fmt.Errorf("Open Error: %s", err)
				}
			}
		} else {
			return fmt.Errorf("Dir Error: %s", err)
		}
	}
	return nil
}
