package cli

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Process(inputs []string, options *Options) error {
	name := nameFile(inputs)

	err := combineMarkdown(inputs, name+".md")
	if err != nil {
		return fmt.Errorf("could not combine markdown files: %w", err)
	}

	docName := name + ".docx"
	fmt.Println(options)
	if options.output != "" {
		docName = options.output
	}

	if err := convertToDoc(name+".md", docName, options); err != nil {
		return fmt.Errorf("could not convert to doc: %w", err)
	}

	if err := fiddleWithDoc(name+".docx", options); err != nil {
		return fmt.Errorf("could not style doc: %w", err)

	}

	if err := deleteMarkdown(name + ".md"); err != nil {
		fmt.Println("\033[33m", fmt.Errorf("could not delete temp markdown file: %w", err))
	}
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
	file, err := os.Create(name)

	if err != nil {
		return fmt.Errorf("could not create tmp.md: %w", err)
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

		_, err := io.Copy(w, titleReader)
		if err != nil {
			return fmt.Errorf("could not write to %v: %w", name, err)
		}

		if _, err := io.Copy(w, reader); err != nil {
			return fmt.Errorf("could not write to %v: %w", name, err)
		}

		w.Flush()

	}

	return nil
}

func convertToDoc(markdownName string, docName string, options *Options) error {
	cmd := exec.Command("pandoc", markdownName, `--from=markdown`, `--to=docx`, `-o `+docName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("couldn't pandoc: %w", err)
	}

	return nil
}

func fiddleWithDoc(name string, options *Options) error {
	return nil
}

func deleteMarkdown(name string) error {
	err := os.Remove(name)
	return err
}
