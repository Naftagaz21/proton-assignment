package htmlutils

import (
	"fmt"
)

var SeparatorMdSeq = []byte("\n\n---")

const startSequenceTemplate = `
<!DOCTYPE html>
<html lang="en"><head>
<meta http-equiv="content-type" content="text/html; charset=UTF-8">
<meta charset="utf-8">
<title>%s</title>
<style>html 
   * { font-family: "Arial", sans-serif; }
   hr { border: solid 1px #ccc; margin-bottom: 50px; margin-top: 50px}
   body { width: 750px; margin-left: auto; margin-right: auto; }
   h1 { text-align: right; color: #6d4aff; margin-bottom: 50px; }
   h2 { text-align: center; color: #372580; margin-bottom: 50px; }
   p { text-align: justify; }
   .pagination-wrapper {display: flex; margin-top: 24px; margin-bottom: 24px;}
   .both {justify-content: space-between;}
   .next {justify-content: flex-end;}
   .previous {justify-content: flex-start;}
   </style>
</head>
<body>
<h1>%s</h1>
`

const endSequenceTemplate = `
<div class="pagination-wrapper %s">
%s
%s
</div>
</body>`

type EndSequenceType int

const (
	NextPrevious EndSequenceType = iota
	Next
	Previous
	None
)

func determineEndType(postLength int, idx int) EndSequenceType {
	if postLength <= 1 {
		return None
	} else if idx == 0 {
		return Next
	} else if idx > 0 && idx < postLength-1 {
		return NextPrevious
	} else {
		return Previous
	}
}

type Htmlutils interface {
	GetStartSequence() string
	GenerateEndSequence(postLength int, pageIdx int) string
	GetFilename(idx int) string
}

type htmlutils struct {
	startSequence string
}

func (util *htmlutils) GetFilename(idx int) string {
	return fmt.Sprintf("page-%d.html", idx)
}

func (util *htmlutils) GetStartSequence() string {
	return util.startSequence
}

func (util *htmlutils) GenerateEndSequence(postLength int, pageIdx int) string {
	endType := determineEndType(postLength, pageIdx)
	prevPage := util.GetFilename(pageIdx - 1)
	nextPage := util.GetFilename(pageIdx + 1)

	var className, linkPrev, linkNext string

	switch endType {
	case NextPrevious:
		className = "both"
		linkPrev = fmt.Sprintf(`<a href="%s">Previous Page</a>`, prevPage)
		linkNext = fmt.Sprintf(`<a href="%s">Next Page</a>`, nextPage)
	case Next:
		className = "next"
		linkNext = fmt.Sprintf(`<a href="%s">Next Page</a>`, nextPage)
	case Previous:
		className = "previous"
		linkPrev = fmt.Sprintf(`<a href="%s">Previous Page</a>`, prevPage)
	}

	return fmt.Sprintf(endSequenceTemplate, className, linkPrev, linkNext)
}

func NewHtmlutils(title string) Htmlutils {
	return &htmlutils{
		startSequence: fmt.Sprintf(startSequenceTemplate, title, title),
	}
}
