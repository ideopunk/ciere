package cli

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	// "github.com/gogap/config"
	// "github.com/gogap/go-pandoc/pandoc"
)

func Process(inputs []string, options *Options) error {
	name := nameFile(inputs)

	err := combineMarkdown(inputs, name)
	if err != nil {
		return fmt.Errorf("could not combine markdown files: %v", err)
	}

	convertToDoc(name, options)
	deleteMarkdown(name)
	return nil
}

func getTitle(fileName string) string {
	s := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	pretitle := strings.Replace(s, "-", "", -1)
	title := "**" + strings.Title(pretitle) + "**\n\n"
	return title
}

func nameFile(inputs []string) string {
	fileName := ""

	for ind, input := range inputs {
		base := filepath.Base(input)
		cleanBase := strings.TrimSuffix(base, filepath.Ext(input))
		if ind != 0 {
			fileName += "_"
		}
		fileName += cleanBase
	}

	return fileName
}

func combineMarkdown(inputs []string, name string) error {
	mdName := name + ".md"
	file, err := os.Create(mdName)

	if err != nil {
		return fmt.Errorf("could not create tmp.md: %v", err)
	}

	defer file.Close()

	for i, input := range inputs {
		title := getTitle(input)
		if i != 0 {
			title = "\n\n" + title
		}
		titleReader := strings.NewReader(title)
		reader := strings.NewReader(input)

		w := bufio.NewWriter(file)

		uh, err := io.Copy(w, titleReader)
		println(uh)
		if err != nil {
			return fmt.Errorf("could not write to %v: %v", mdName, err)
		}

		if _, err := io.Copy(w, reader); err != nil {
			return fmt.Errorf("could not write to %v: %v", mdName, err)
		}

		w.Flush()

	}

	return nil
}

func convertToDoc(name string, options *Options) error {
	docName := name + ".docx"
	println(docName)
	// fetch := pandoc.FetcherOptions{}
	// conv := pandoc.ConvertOptions{From: "markdown", To: "docx"}
	// pdoc, err := pandoc.New(&config.Config{})

	// if err != nil {
	// 	return fmt.Errorf("couldn't pandoc: %v", err)
	// }

	// pdoc.Convert(fetch, conv)

	return nil
}

func deleteMarkdown(name string) {
	mdName := name + ".md"
	println(mdName)
}
