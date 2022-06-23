package cli

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/nguyenthenguyen/docx"
)

func Process(inputs []string, options *Options) error {
	name := nameFile(inputs)
	markdownName := name + ".md"
	docName := name + ".docx"

	err := combineMarkdown(inputs, name+".md", options.double)
	if err != nil {
		return fmt.Errorf("could not combine markdown files: %w", err)
	}

	fmt.Println(options)
	if options.output != "" {
		docName = options.output
	}

	if err := convertToDoc(markdownName, docName, options); err != nil {
		return fmt.Errorf("could not convert to doc: %w", err)
	}

	if err := fiddleWithHeader(docName, options); err != nil {
		return fmt.Errorf("could not style doc: %w", err)

	}

	if err := deleteMarkdown(markdownName); err != nil {
		fmt.Println("\033[33m", fmt.Errorf("could not delete temp markdown file: %w", err))
	}
	return nil
}

func getTitle(fileName string, ind int) string {
	s := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	pretitle := strings.Replace(s, "-", "", -1)
	title := "**" + strings.Title(pretitle) + "**\n\n"
	if ind != 0 {
		// for pieces after the first one, add a pagebreak
		title = "\n```{=openxml}\n<w:p><w:r><w:br w:type=\"page\"/></w:r></w:p>\n```\n" + title
	}
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

func fiddleWithHeader(name string, options *Options) error {
	r, err := docx.ReadDocxFile(name)
	if err != nil {
		return fmt.Errorf("couldn't open docx we just made: %w", err)
	}

	docx := r.Editable()


	docx.Replace("1", "ONE!!!", 1)
	docx.ReplaceFooter("", "FOOT")
	
	// by default my pandoc reference doc includes page numbers. If `options.page` is false, we go in and remove them
	if !options.page {
		for i := 0; i < 10; i++ {
			err := docx.ReplaceHeader(fmt.Sprint(i), "")
			if err != nil {
				return fmt.Errorf("could not update header: %w", err)
			}
		}

	}

	if options.author {
		docx.ReplaceHeader(" ", "Conor Barnes")
	}

	r.Close()
	return nil
}

func deleteMarkdown(name string) error {
	err := os.Remove(name)
	return err
}
