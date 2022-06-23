package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func Process(inputs []string, options *Options) error {
	name, err := nameFile(inputs, options.output)
	if err != nil {
		return fmt.Errorf("could not name file: %w", err)
	}

	markdownName := strings.TrimSuffix(name, ".docx") + ".md"

	if err := combineMarkdown(inputs, markdownName, options.double); err != nil {
		return fmt.Errorf("could not combine markdown files: %w", err)
	}

	if err := convertToDoc(markdownName, name, options); err != nil {
		return fmt.Errorf("could not convert to doc: %w", err)
	}

	if err := deleteMarkdown(markdownName); err != nil {
		fmt.Println("\033[33m", fmt.Errorf("could not delete temp markdown file: %w", err))
	}

	success(name)
	return nil
}

// get the title of a piece for inclusion inside the final document. also include xml for a pagebreak if need be.
func getTitle(fileName string, ind int) string {
	s := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	pretitle := strings.Replace(s, "-", "", -1)
	caser := cases.Title(language.English)
	title := "**" + caser.String(pretitle) + "**\n\n"
	if ind != 0 {
		// for pieces after the first one, add a pagebreak
		title = "\n```{=openxml}\n<w:p><w:r><w:br w:type=\"page\"/></w:r></w:p>\n```\n" + title
	}
	return title
}

// create the name that will be used for the files.
func nameFile(inputs []string, userAssignedOutput string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not access user home dir: %w", err)
	}

	fileName := ""

	if userAssignedOutput != "" {
		fileName = homeDir + "/Documents/Writing/submissions/" + userAssignedOutput

	} else {
		for ind, input := range inputs {
			base := filepath.Base(input)
			cleanBase := strings.TrimSuffix(base, filepath.Ext(input))
			if ind != 0 {
				fileName += "_"
			}
			fileName += cleanBase
		}
		fileName = homeDir + "/Documents/Writing/submissions/" + fileName + ".docx"

	}

	return fileName, nil
}

func combineMarkdown(inputs []string, name string, double bool) error {
	file, err := os.Create(name)

	if err != nil {
		return fmt.Errorf("could not create %v: %w", name, err)
	}

	defer file.Close()

	w := bufio.NewWriter(file)

	if double {
		startReader := strings.NewReader("::: {custom-style=\"Double\"}\n")
		if _, err := io.Copy(w, startReader); err != nil {
			return fmt.Errorf("could not write to %v: %w", name, err)
		}
	}

	for i, input := range inputs {
		title := getTitle(input, i)

		titleReader := strings.NewReader(title)

		file, err := os.Open(input)
		if err != nil {
			return fmt.Errorf("could not open %v: %w", input, err)
		}

		defer file.Close()

		reader := bufio.NewReader(file)

		if _, err := io.Copy(w, titleReader); err != nil {
			return fmt.Errorf("could not write to %v: %w", name, err)
		}

		if _, err := io.Copy(w, reader); err != nil {
			return fmt.Errorf("could not write to %v: %w", name, err)
		}

	}

	if double {
		endReader := strings.NewReader("\n:::")
		if _, err := io.Copy(w, endReader); err != nil {
			return fmt.Errorf("could not write to %v: %w", name, err)
		}
	}

	w.Flush()

	return nil
}

func convertToDoc(markdownName string, docName string, options *Options) error {
	cmd := exec.Command("pandoc", markdownName,
		`--from=markdown+backtick_code_blocks+raw_attribute+hard_line_breaks`,
		`--to=docx`,
		`--output=`+docName,
		`--reference-doc=reference.docx`,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("couldn't pandoc: %w", err)
	}

	return nil
}

func deleteMarkdown(name string) error {
	err := os.Remove(name)
	return err
}

func success(path string) {
	base := filepath.Base(path)
	fmt.Printf("\033[32m%s has been succesfully created!\n", base)
}
