package main

import (
    "bufio"
    "fmt"
    "net/http"
    "regexp"
    "strings"
    "encoding/json"
)

// Define token types
type TokenType int

const (
	Header TokenType = iota
	Paragraph
	LineBreak
	Bold
	Italic
	Link
	ListItem
	Text
)

// Token struct
type Token struct {
	Type  TokenType
	Value string
}

// Tokenize markdown input
func tokenize(input string) []Token {
	tokens := []Token{}
	scanner := bufio.NewScanner(strings.NewReader(input))

	for scanner.Scan() {
		line := scanner.Text()

		// Detect empty line for line breaks
		if line == "" {
			tokens = append(tokens, Token{Type: LineBreak, Value: "<br/>"})
			continue
		}

		// Detect headers and remove Markdown syntax
		headerRegex := regexp.MustCompile(`^(#+)\s+(.*)`)
		if headerMatch := headerRegex.FindStringSubmatch(line); headerMatch != nil {
			tokens = append(tokens, Token{Type: Header, Value: headerMatch[2]})
			continue
		}

		// Detect list items and remove Markdown syntax
		listItemRegex := regexp.MustCompile(`^\*\s+(.*)`)
		if listItemMatch := listItemRegex.FindStringSubmatch(line); listItemMatch != nil {
			tokens = append(tokens, Token{Type: ListItem, Value: listItemMatch[1]})
			continue
		}

		// Parse inline elements (bold, italic, link) along with text
		parseInlineElements(&tokens, line)
	}

	return tokens
}

func parseInlineElements(tokens *[]Token, line string) {
	// Regex for bold, italic, and links
	boldRegex := regexp.MustCompile(`\*\*(.*?)\*\*`)
	italicRegex := regexp.MustCompile(`\*(.*?)\*`)
	linkRegex := regexp.MustCompile(`\[(.*?)\]\((.*?)\)`)

	// Parse line for links, bold, and italic text
	for {
		linkMatch := linkRegex.FindStringSubmatchIndex(line)
		boldMatch := boldRegex.FindStringSubmatchIndex(line)
		italicMatch := italicRegex.FindStringSubmatchIndex(line)

		// Find the earliest match
		lowest := len(line)
		matchType := "none"
		matchIndex := []int{}

		if linkMatch != nil && linkMatch[0] < lowest {
			lowest = linkMatch[0]
			matchType = "link"
			matchIndex = linkMatch
		}
		if boldMatch != nil && boldMatch[0] < lowest {
			lowest = boldMatch[0]
			matchType = "bold"
			matchIndex = boldMatch
		}
		if italicMatch != nil && italicMatch[0] < lowest {
			lowest = italicMatch[0]
			matchType = "italic"
			matchIndex = italicMatch
		}

		if matchType == "none" {
			// No more matches, add the rest of the line as text
			if len(strings.TrimSpace(line)) > 0 {
				*tokens = append(*tokens, Token{Type: Text, Value: line})
			}
			break
		}

		// Add text before match as a paragraph token if it's not empty
		beforeMatch := line[:lowest]
		if len(strings.TrimSpace(beforeMatch)) > 0 {
			*tokens = append(*tokens, Token{Type: Text, Value: beforeMatch})
		}

		// Add matched token
		switch matchType {
		case "link":
			*tokens = append(*tokens, Token{Type: Link, Value: fmt.Sprintf("Text: %s, URL: %s", line[matchIndex[2]:matchIndex[3]], line[matchIndex[4]:matchIndex[5]])})
		case "bold":
			*tokens = append(*tokens, Token{Type: Bold, Value: line[matchIndex[2]:matchIndex[3]]})
		case "italic":
			*tokens = append(*tokens, Token{Type: Italic, Value: line[matchIndex[2]:matchIndex[3]]})
		}

		// Update line to continue after the match
		line = line[matchIndex[1]:]
	}
}

// Markdown to HTML parser
func parseToHTML(tokens []Token) string {
	var html strings.Builder
	listOpen := false // Track if we are inside a list

	for _, token := range tokens {
		switch token.Type {
		case Header:
			level := strings.Count(token.Value, "#") // Determine header level by count of '#'
			content := strings.TrimSpace(token.Value[len("#"):])
			html.WriteString(fmt.Sprintf("<h%d>%s</h%d>\n", level, content, level))
		case Paragraph:
			html.WriteString(fmt.Sprintf("<p>%s</p>\n", token.Value))
		case LineBreak:
			html.WriteString("<br/>\n")
		case Bold:
			html.WriteString(fmt.Sprintf("<strong>%s</strong>", token.Value))
		case Italic:
			html.WriteString(fmt.Sprintf("<em>%s</em>", token.Value))
		case Link:
			parts := strings.SplitN(token.Value, ", URL: ", 2)
			text := strings.TrimPrefix(parts[0], "Text: ")
			url := parts[1]
			html.WriteString(fmt.Sprintf("<a href=\"%s\">%s</a>", url, text))
		case ListItem:
			if !listOpen {
				html.WriteString("<ul>\n")
				listOpen = true
			}
			html.WriteString(fmt.Sprintf("<li>%s</li>\n", token.Value))
		case Text:
			html.WriteString(fmt.Sprintf("%s", token.Value))
		}
	}

	if listOpen {
		html.WriteString("</ul>\n") // Close the list if it was opened
	}

	return html.String()
}

func markdownToHTMLHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var data struct {
        Markdown string `json:"markdown"`
    }
    err := json.NewDecoder(r.Body).Decode(&data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Process the Markdown
    tokens := tokenize(data.Markdown)
    html := parseToHTML(tokens)

    // Return the HTML
    fmt.Fprintf(w, html)
}

func main() {
    http.HandleFunc("/convert", markdownToHTMLHandler)
    fmt.Println("Server is listening on port 8080...")
    http.ListenAndServe(":8080", nil)
}


