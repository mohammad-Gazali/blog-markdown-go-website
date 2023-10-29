package server

import (
	"net/http"
	"strings"
	"fmt"
	"html"
	"regexp"
)

// ============================ Web Utils ============================

func InternalServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	RenderTemplate(w, RenderContext{"Title": "Server Error"}, "500.html")
}

func NotFoundError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	RenderTemplate(w, RenderContext{"Title": "Not Found"}, "404.html")
}


// ============================ Markdown Utils ============================

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
	pureString := html.EscapeString(strings.Trim(s, " " + "\r" + "\t"))

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

	} else if pureString[0] == '>' && !isContent {
		// blockquote tag case
		return WrapWithTag(MarkdownToHTML(pureString[1:], true, false), "blockquote")

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
		s = reInlineCode.ReplaceAllString(s, WrapWithTag("$1", "code"))
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