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
	Image        string
	Tag          string
	Comment      string
	HowManyTimes int64
}

func render(c Container) {
	type Itercontainer struct {
		Iter int64
		Container
	}
	text := `---
apiVersion: v1
kind: ConfigMap
metadata:
  name: container-cm-{{.Iter}}
data: 
  image: {{.Image}}:{{.Tag}}
  comment: |
    {{.Comment}}

`

	t, err := template.New("text").Parse(text)
	if err != nil {
		panic(err)
	}

	var i int64
	for i = 0; i < c.HowManyTimes; i++ {
		iterc := Itercontainer{
			Iter: i,
			Container: Container{
				Image:   c.Image,
				Tag:     c.Tag,
				Comment: c.Comment,
			},
		}

		t.Execute(os.Stdout, iterc)

	}

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
		Image:        tree.Get("container.image").(string),
		Tag:          tree.Get("container.tag").(string),
		Comment:      tree.Get("container.comment").(string),
		HowManyTimes: tree.Get("container.how_many_times").(int64),
	}

	render(cont)

}
