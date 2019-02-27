package main

import (
	"os"
	"log"
	"io"
	"net/http"
	"strings"
	"bytes"
	"io/ioutil"
	"time"
	"path"
	"path/filepath"
	"fmt"
	"os/exec"
)

func mkWorkDir() (string, string, error) {
	workdir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln("Cannot get user home")
		return "", "", err
	}
	workdir += "/.local/share/unsplash/"
	olddir := workdir + "old/"
	err = os.MkdirAll(olddir, os.ModePerm)
	if err != nil {
		log.Fatalln("cannot mkdir ", olddir)
		return "", "", err
	}
	return workdir, olddir, nil
}

func backupOldFile(wd string, od string) {
	Depth := strings.Count(wd, "/")
	_ = filepath.Walk(wd, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		depth := strings.Count(p, "/")
		if Depth != depth {
			return nil
		}
		if filepath.Ext(p) == ".jpg" {
			os.Rename(p, path.Join(od, info.Name()))
		}
		return nil
	})
}

func getImg(url string, fn string) error {
	out , err := os.Create(fn)
	if err != nil {
		return err
	}
	defer out.Close()
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	pix, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	_, err = io.Copy(out, bytes.NewReader(pix))
	return nil
}

func main() {
	workdir, olddir, err := mkWorkDir()
	if err != nil {
		log.Println("err:", err)
	}
	backupOldFile(workdir, olddir)
	newfilename := workdir + fmt.Sprintf("%d", time.Now().Unix()) + ".jpg"
	getImg("https://source.unsplash.com/1280x1024", newfilename)
	exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-uri", "file://" + newfilename).Run()
}
