package main

//https://colobu.com/2021/01/17/go-embed-tutorial/

import (
	"embed"
	_ "embed"
	"fmt"
	"io/fs"
	"io/ioutil"
)

//go:embed assets/hello.txt
var es string

//go:embed assets/hello.txt
var es1 []byte

//go:embed assets/*
//go:embed assets/hello.txt assets/hello1.txt
//go:embed assets/logs/hello.log
var ef embed.FS

//go:embed assets/*
var ef1 embed.FS

func main() {
	fmt.Printf("嵌入为字符串：%s\n", es)
	fmt.Printf("嵌入为字节数组：%s\n", string(es1))

	efs, _ := ef.ReadFile("assets/logs/hello.log")
	fmt.Printf("嵌入为文件：%s\n", string(efs))

	fmt.Printf("嵌入为目录文件：%s\n", ef)
	dirs, _ := ef.ReadDir("assets")
	for _, d := range dirs {
		fmt.Printf("Read dir file: %s\n", d.Name())
	}

	fmt.Println("返回子文件夹作为新的文件系统")
	ps, _ := fs.Sub(ef, "assets")
	hi, _ := ps.Open("logs/hello.log")
	data, _ := ioutil.ReadAll(hi)
	fmt.Println(string(data))
}
