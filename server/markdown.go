package server

import (
	"fmt"
	"html"
	"html/template"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type MarkdownFile struct {
	Slug    string
	Content string
}

type MarkdownFileInfo struct {
	Title       string
	Description string
	Slug        string
	ImageUrl    string
	CreatedAt   time.Time
	Html        template.HTML
}

func (mf MarkdownFile) Parse() *MarkdownFileInfo {
	infoLines := strings.Split(mf.Content, "\n")

	result := MarkdownFileInfo{}

	result.Slug = mf.Slug

	infoSepCount := 0
	var startHtmlIndex int

	for i, line := range infoLines {
		if strings.Contains(line, "---") {
			infoSepCount++

		} else if strings.Contains(line, "title") {
			result.Title = GetLineInfo(line, "title")			
			
		} else if strings.Contains(line, "description") {
			result.Description = GetLineInfo(line, "description")

		} else if strings.Contains(line, "created_at") {
			if dateInfo := strings.Split(GetLineInfo(line, "created_at"), "/"); len(dateInfo) >= 3 {
				
				year, err := strconv.ParseUint(strings.Trim(dateInfo[0], " " + "\r"), 10, 16)
				if err != nil {
					continue
				}

				month, err := strconv.ParseUint(strings.Trim(dateInfo[1], " " + "\r"), 10, 16)
				if err != nil {
					continue
				}

				day, err := strconv.ParseUint(strings.Trim(dateInfo[2], " " + "\r"), 10, 16)
				if err != nil {
					continue
				}

				result.CreatedAt = time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, time.UTC)
			}

		} else if strings.Contains(line, "image") {
			result.ImageUrl = GetLineInfo(line, "image")
		}

		if infoSepCount >= 2 {
			startHtmlIndex = i + 1
			break
		}
	}

	var htmlResult string
	var isCodeBlock bool
	var currentCodeBlockContent string
	var isOrderdList bool
	var isUnorderdList bool

	for _, line := range infoLines[startHtmlIndex:] {
		if strings.Trim(line, " " + "\r" + "\t") == "```" {
			if isCodeBlock {
				currentCodeBlockContent += "</pre></code>"
				htmlResult += currentCodeBlockContent
				currentCodeBlockContent = ""
				continue
				
			} else {
				currentCodeBlockContent += "<code><pre>"
			}
			isCodeBlock = !isCodeBlock
		}

		// handling lists
		reOrderedListItem := regexp.MustCompile(`(\d)\.\s+(.*)`)
		reUnOrderedListItem := regexp.MustCompile(`-\s+(.*)`)

		if reOrderedListItem.MatchString(line) {
			if !isOrderdList {
				orderedListCounter := reOrderedListItem.ReplaceAllString(line, "$1")
				htmlResult += fmt.Sprintf(`<ol start="%s">`, orderedListCounter)
			}
			isOrderdList = true

		} else if isOrderdList {
			isOrderdList = false
			htmlResult += "</ol>"
		}

		if reUnOrderedListItem.MatchString(line) {
			if !isUnorderdList {
				htmlResult += "<ul>"
			}

			isUnorderdList = true
		} else if isUnorderdList {
			isUnorderdList = false
			htmlResult += "</ul>"
		}

		if isCodeBlock {
			if strings.Trim(line, " " + "\r" + "\t") != "```" {
				currentCodeBlockContent += MarkdownToHTML(line, true, false)
			}
		} else if isOrderdList {
			htmlResult += WrapWithTag(reOrderedListItem.ReplaceAllString(PureString(line), "$2"), "li")

		} else if isUnorderdList {
			htmlResult += WrapWithTag(reUnOrderedListItem.ReplaceAllString(PureString(line), "$1"), "li")

		} else {
			htmlResult += MarkdownToHTML(line, false, true)
		}
	}

	// ending code tag if it is not ended in the markdown file
	if currentCodeBlockContent != "" {
		htmlResult += currentCodeBlockContent + "</pre></code>"
	}

	// ending ordered list ....
	if isOrderdList {
		htmlResult += "</ol>"
	}

	// ending unordered list ....
	if isUnorderdList {
		htmlResult += "</ul>"
	}


	result.Html = template.HTML(htmlResult)

	return &result
}

func GetAllMarkdownFiles() []*MarkdownFileInfo {
	files, err := os.ReadDir("markdown")

	if err != nil {
		return nil
	}

	var result []*MarkdownFileInfo

	for _, f := range files {

		content, err := os.ReadFile("markdown/" + f.Name())

		if err != nil {
			return nil
		}

		result = append(result, MarkdownFile{
			Slug:    strings.Split(f.Name(), ".")[0],
			Content: string(content),
		}.Parse())
	}

	return result
}

func GetMarkdownBySlug(slug string) *MarkdownFileInfo {
	content, err := os.ReadFile("markdown/" + slug + ".md")

	if err != nil {
		return nil
	}

	return MarkdownFile{
		Slug:    slug,
		Content: string(content),
	}.Parse()
}


// ============ utils ============

// this function return the info in string line
//
// for example:
//
// - GetLineInfo("name: Mohammed Algazali", "name")    >> "Mohammed Algazali"
//
// - GetLineInfo("created_at: 2023/4/1", "created_at") >> "2023/4/1"
//
// - GetLineInfo("job: programmer", "jop")             >> ""  # because we failed in key name writing
// 
func GetLineInfo(line, key string) string {
	if info := strings.Split(line, key); len(info) >= 2 {
		// here we remove ":", spaces and "\r" characters
		return strings.Trim(info[1], " " + ":" + "\r")
	} else {
		return ""
	}
}

func MarkdownToHTML(s string, isContent, isParagraph bool) string {
	pureString := html.EscapeString(PureString(s))

	if pureString == "" {
		return ""
	}

	if pureString == "---" && !isContent {
		// horizontal line case
		return "<hr />"

	} else if pureString[0] == '#' && !isContent {
		// header tag case
		num := FirstSequenceCharsCount(pureString, '#')
		
		if num > 6 {
			num = 6
		}

		if len(pureString) > num + 1 && pureString[num] == ' ' {
			return WrapWithTag(MarkdownToHTML(pureString[num:], true, false), fmt.Sprintf("h%v", num))
		} else {
			// here we ecscape the header case in the when there is no space after the valid number of "#" characters
			// so we added to the function isContent as true to continue to the next case after the header because it is the last else if that cares about isContent var
			//! make sure you don't move this case up, keep it last case that cares about isContent
			return MarkdownToHTML(s, true, isParagraph)
		}

	} else if PureString(s)[0] == '>' && !isContent {
		// blockquote tag case
		return WrapWithTag(MarkdownToHTML(pureString[5:], true, false), "blockquote")

	} else {
		// inline tags like bold and italic and others...
		reBold := regexp.MustCompile(`(\*{2}([^\*\n]*)\*{2})|(_{2}([^_\n]*)_{2})`)
		reItalic := regexp.MustCompile(`(\*([^\*\n]+)\*)|(_([^_\n]+)_)`)
		reInlineCode := regexp.MustCompile("`(.+)`")
		reImage := regexp.MustCompile(`!\[(.*)\]\((.*)\)`)
		reLink := regexp.MustCompile(`[^!]?\[(.*)\]\((.*)\)`)
		reStrikeThrough := regexp.MustCompile(`~{2}(.*?)~{2}`)
		reSub := regexp.MustCompile(`~(.*?)~`)
		reSup := regexp.MustCompile(`\^(.*?)\^`)
		reMark := regexp.MustCompile(`==(.*?)\==`)

		//! we must start with bold because if we don't parse it first it will be considered as italic
		s = reBold.ReplaceAllString(s, WrapWithTag("$2$4", "strong"))
		s = reItalic.ReplaceAllString(s, WrapWithTag("$2$4", "em"))
		s = reInlineCode.ReplaceAllString(s, WrapWithTag(WrapWithTag("$1", "code"), "pre"))
		s = reImage.ReplaceAllString(s, `<img src="$2" alt="$1">`)
		s = reLink.ReplaceAllString(s, `<a href="$2">$1</a>`)
		//! we must start with strike through because if we don't parse it first it will be considered as sub
		s = reStrikeThrough.ReplaceAllString(s, WrapWithTag("$1", "del"))
		s = reSub.ReplaceAllString(s, WrapWithTag("$1", "sub"))
		s = reSup.ReplaceAllString(s, WrapWithTag("$1", "sup"))
		s = reMark.ReplaceAllString(s, WrapWithTag("$1", "mark"))

		if isParagraph {
			return WrapWithTag(s, "p")
		}

		return s
	}
}

func WrapWithTag(content, tag string) string {
	return fmt.Sprintf("<%s>%s</%s>", tag, content, tag)
}

func FirstSequenceCharsCount(s string, c rune) int {
	var result int

	for _, char := range s {
		if result != 0 && char != c {
			return result
		} else if char == c {
			result++
		}
	}

	return result
}

func PureString(s string) string {
	return strings.Trim(s, " " + "\r" + "\t")
}