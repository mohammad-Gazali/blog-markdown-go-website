package server

import (
	"os"
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
	Html        string
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
		}
	}

	// TODO: handle parsing html
	result.Html = strings.Join(infoLines[startHtmlIndex:], "\n")

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
