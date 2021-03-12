package main

import (
	// "fmt"
	"os"
	"text/template"

	// "github.com/pelletier/go-toml"
	"io/fs"
	"io/ioutil"
	"strings"

	"github.com/pelletier/go-toml"
)

type Container struct {
	Image   string
	Tag     string
	Comment string
}

func render(c Container) {
	text := `apiVersion: v1
kind: ConfigMap
metadata:
  name: container-cm
data: 
  image: {{.Image}}:{{.Tag}}
  comment:|
    {{.Comment}}
`
	t, err := template.New("text").Parse(text)
	if err != nil {
		panic(err)
	}

	t.Execute(os.Stdout, c)

}

func main() {
	// fmt.Println("vim-go")
	files, _ := ioutil.ReadDir(".")
	tomlFiles := make([]fs.FileInfo, 0)
	for _, f := range files {
		if strings.HasSuffix(f.Name(), "toml") {
			tomlFiles = append(tomlFiles, f)
		}
	}
	if len(tomlFiles) != 1 {
		panic("There should be one file ending with the toml extension")
	}

	// reading the file
	data, err := ioutil.ReadFile(tomlFiles[0].Name())
	if err != nil {
		panic("could not open file")
	}

	tree, err := toml.Load(string(data))
	if err != nil {
		panic(err)
	}

	cont := Container{
		Image:   tree.Get("container.image").(string),
		Tag:     tree.Get("container.tag").(string),
		Comment: tree.Get("container.comment").(string),
	}

	render(cont)

}
