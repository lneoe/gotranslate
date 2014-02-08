package history

// package main

import (
    "fmt"
    "os"
    "os/user"
    "strings"
)

func findHomePath() string {
    u, _ := user.Current()
    homePath := u.HomeDir
    a := []string{homePath, ".translate"}
    path := strings.Join(a, "/")
    //检查目录 .translate 是否存在
    if !exists(path) {
        // fmt.Println(".translate 目录不存在，开始创建..")
        os.Mkdir(path, 0766)
    }
    return path
}

func exists(path string) bool {
    _, err := os.Stat(path)
    // return err == nil || os.IsExist(err) || !os.IsNotExist(err)
    return err == nil || os.IsExist(err)
}

type HistoryFile struct {
    file *os.File
    // path string
}

func (h *HistoryFile) open() {
    homePath := findHomePath()
    path := homePath + "/history"
    // a := []string{homePath, "history"}
    // path := strings.Join(a, "/")
    file, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        panic("openfile error")
    } else {
        h.file = file
    }
}

func (h *HistoryFile) close() {
    h.file.Close()
}

func (h *HistoryFile) Add(word string) {
    fmt.Println("add history")
    // file := *h.file
    h.open()
    h.file.WriteString(word)
    h.close()

}

func (h *HistoryFile) Clear() {
    h.open()
    h.file.Truncate(0)
    h.close()

}

var (
    History HistoryFile
)
