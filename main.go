package table

import (
	"bytes"
	"strings"

	"github.com/fatih/color"
)

var MaxCellWidth = 40
var LineEndTag = "+"
var LineBody = "-"
var LineDivider = "|"

func trimTableCell(headers []string, content [][]string) {
	for i, header := range headers {
		headers[i] = strings.TrimSpace(header)
	}
	for i, c := range content {
		for j, s := range c {
			content[i][j] = strings.TrimSpace(s)
		}
	}
}

func Show(headers []string, content [][]string) {
	trimTableCell(headers, content)
	// 寻找每列最长字符
	maxWidths := make([]int, len(headers), len(headers))
	for i, header := range headers {
		maxWidths[i] = len(header)
	}
	for _, row := range content {
		for i, col := range row {
			if i < len(maxWidths) {
				max := maxWidths[i]
				current := len(col)
				if current > max {
					maxWidths[i] = current
				}
			}
		}
	}
	for i, width := range maxWidths {
		if width > MaxCellWidth {
			maxWidths[i] = MaxCellWidth
		}
	}

	headerWriter := color.New(color.FgBlue, color.Bold)
	// 输出header
	headerWriter.Println(tableBorder(maxWidths))

	headerWriter.Printf("| %-*s |", maxWidths[0]-findNonAsciiCount(headers[0]), headers[0])
	for i := 1; i < len(headers); i++ {
		headerWriter.Printf(" %-*s |", maxWidths[i]-findNonAsciiCount(headers[i]), headers[i])
	}

	if len(headers) != 0 {
		headerWriter.Println(tableBorder(maxWidths))
	}
	contentWriter := color.New(color.FgBlue)
	for _, c := range content {
		contentWriter.Printf(row(c, maxWidths))
	}

	headerWriter.Println(tableBorder(maxWidths))
}

func tableBorder(maxWidths []int) string {
	bufferString := bytes.NewBufferString("")
	for i, width := range maxWidths {
		width += 2
		if i == len(maxWidths)-1 {
			bufferString.WriteString(borderWithEnd(width, true))
		} else {
			bufferString.WriteString(borderWithEnd(width, false))
		}
	}
	return bufferString.String()
}

func row(row []string, maxWidth []int) string {
	maxLine := 0
	for i, cell := range row {
		if i >= len(maxWidth) {
			break
		}
		lines := howManyLine(cell, maxWidth[i])
		if lines > maxLine {
			maxLine = lines
		}
	}
	bufferString := bytes.NewBufferString("")
	for i := 0; i < maxLine; i++ {
		bufferString.WriteString(LineDivider)
		for j, s := range row {
			if j < len(maxWidth) {
				bufferString.WriteString(" ")
				bufferString.WriteString(fixStr(s, i, maxWidth[j]))
				bufferString.WriteString(" ")
				bufferString.WriteString(LineDivider)
			}
		}
		bufferString.WriteString("\n")
	}
	return bufferString.String()
}

func borderWithEnd(length int, end bool) string {
	bufferString := bytes.NewBufferString(LineEndTag)
	for i := 0; i < length; i++ {
		bufferString.WriteString(LineBody)
	}
	if end {
		bufferString.WriteString(LineEndTag)
	}
	return bufferString.String()
}

func fixStr(str string, line int, size int) string {
	currentLine := line * size
	overLine := (line + 1) * size
	if len(str) < currentLine {
		bufferString := bytes.NewBufferString("")
		for i := 0; i < size; i++ {
			bufferString.WriteByte(' ')
		}
		return bufferString.String()
	} else if len(str) > overLine { // 可以直接截断
		return string([]rune(str)[currentLine:overLine])
	} else {
		endLength := len([]rune(str))
		content := string([]rune(str)[currentLine:endLength])
		bufferString := bytes.NewBufferString(content)
		for i := 0; i < overLine-endLength; i++ {
			bufferString.WriteByte(' ')
		}
		return bufferString.String()
	}
}

func howManyLine(str string, cut int) int {
	line := len(str) / (cut + 1)
	return line + 1
}

func findNonAsciiCount(str string) int {
	count := 0
	for _, r := range str {
		if r > 255 {
			count++
		}
	}
	return count
}
