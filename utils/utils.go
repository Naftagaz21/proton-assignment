package utils

import (
	"bytes"
	"math"
	"regexp"
	"sort"
	"time"
)

const timeFormat = `2006-01-02`

var publicationDateRegexPattern *regexp.Regexp = regexp.MustCompile(`Published on (\d{4}-\d{2}-\d{2})`)

type FormatUtils interface {
	SortByPublicationDate(mdData [][]byte)
	SplitByArticleLimit(data [][]byte, limit int) [][]byte
}

type formatUtils struct{}

func extractPublicationDate(data []byte) string {
	match := publicationDateRegexPattern.FindSubmatch(data)
	if len(match) < 2 {
		return ""
	}
	return string(match[1])
}

func (*formatUtils) SortByPublicationDate(mdData [][]byte) {
	sort.Slice(mdData, func(i, j int) bool {
		date1 := extractPublicationDate(mdData[i])
		date2 := extractPublicationDate(mdData[j])

		time1, err := time.Parse(timeFormat, date1)
		if err != nil {
			return false
		}

		time2, err := time.Parse(timeFormat, date2)
		if err != nil {
			return false
		}

		return time1.After(time2)
	})
}

func (*formatUtils) SplitByArticleLimit(data [][]byte, limit int) [][]byte {
	var pages [][]byte

	// A bit hacky
	if limit == 0 {
		limit = math.MaxInt
	}

	for i := 0; i < len(data); i += limit {
		end := i + limit
		if end > len(data) {
			end = len(data)
		}
		page := bytes.Join(data[i:end], []byte("\n"))
		pages = append(pages, page)
	}

	return pages
}

func NewFormatUtils() FormatUtils {
	return &formatUtils{}
}
