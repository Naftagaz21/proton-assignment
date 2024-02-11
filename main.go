package main

import (
	"bytes"
	"fmt"
	"gen-blog/htmlutils"
	"gen-blog/iovalidator"
	"gen-blog/utils"
	"os"
	"path/filepath"
	"sync"

	"github.com/alexflint/go-arg"
	"github.com/yuin/goldmark"
)

type GenerateCmd struct {
	Input               string `arg:"required" help:"a valid input folder path"`
	Output              string `arg:"required" help:"a valid output folder path"`
	Title               string `arg:"required" help:"title for the generated blog"`
	PostsPerPage        int    `arg:"--posts-per-page" default:"0" help:"the amount of posts per single page"`
	DisableSeparator    bool   `arg:"--separator, --s" default:"false" help:"disable the separator at the end of each blog post"`
	MultithreadedOutput bool   `arg:"--multithread, --m" default:"false" help:"enable multithreaded output"`
}

type args struct {
	Command *GenerateCmd `arg:"subcommand:generate" `
}

func (args) Description() string {
	return "converts markdown files to html files"
}

func printError(err error) {
	fmt.Println("Error occured:", err)
}

func main() {
	args := args{}
	p := arg.MustParse(&args)
	if p.Subcommand() == nil {
		p.Fail("Missing command")
	}

	validator := iovalidator.NewValidator()
	mdFilepaths, err := validator.ValidateAndProcessInput(args.Command.Input)
	if err != nil {
		printError(err)
		return
	}

	err = validator.ValidateOutput(args.Command.Output)
	if err != nil {
		printError(err)
		return
	}

	var mdBytes [][]byte
	for _, el := range mdFilepaths {
		mdData, err := os.ReadFile(el)
		if err != nil {
			printError(err)
			return
		}
		if !args.Command.DisableSeparator {
			mdData = append(mdData, htmlutils.SeparatorMdSeq...)
		}
		mdBytes = append(mdBytes, mdData)
	}

	formatUtil := utils.NewFormatUtils()
	formatUtil.SortByPublicationDate(mdBytes)
	mdPages := formatUtil.SplitByArticleLimit(mdBytes, args.Command.PostsPerPage)

	htmlUtil := htmlutils.NewHtmlutils(args.Command.Title)

	if args.Command.MultithreadedOutput {
		multithreadedOutput(mdPages, htmlUtil, args)
	} else {
		singleThreadedOutput(mdPages, htmlUtil, args)
	}
}

func singleThreadedOutput(mdPages [][]byte, htmlUtil htmlutils.Htmlutils, args args) {
	var buf bytes.Buffer
	for idx, el := range mdPages {
		buf.WriteString(htmlUtil.GetStartSequence())

		err := goldmark.Convert(el, &buf)
		if err != nil {
			printError(err)
			os.Exit(1)
		}

		filePath := filepath.Join(args.Command.Output, fmt.Sprintf("page-%d.html", idx))
		buf.WriteString(htmlUtil.GenerateEndSequence(len(mdPages), idx))
		err = os.WriteFile(filePath, buf.Bytes(), 0644)
		if err != nil {
			printError(err)
			os.Exit(1)
		}

		buf.Reset()
	}
}

func multithreadedOutput(mdPages [][]byte, htmlUtil htmlutils.Htmlutils, args args) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	errorCh := make(chan error, len(mdPages))

	for idx, el := range mdPages {
		wg.Add(1)
		go func(idx int, el []byte) {
			defer wg.Done()
			var buf bytes.Buffer

			buf.WriteString(htmlUtil.GetStartSequence())

			err := goldmark.Convert(el, &buf)
			if err != nil {
				errorCh <- err
				return
			}

			filePath := filepath.Join(args.Command.Output, fmt.Sprintf("page-%d.html", idx))
			mu.Lock()
			defer mu.Unlock()

			buf.WriteString(htmlUtil.GenerateEndSequence(len(mdPages), idx))
			err = os.WriteFile(filePath, buf.Bytes(), 0644)
			if err != nil {
				errorCh <- err
				return
			}

			buf.Reset()
		}(idx, el)
	}

	go func() {
		wg.Wait()
		close(errorCh)
	}()

	for err := range errorCh {
		printError(err)
		os.Exit(1)
	}
}
