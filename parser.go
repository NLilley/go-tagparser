package parser

import (
	"fmt"
	"unicode"
)

var TextTagName string = "<text>"
var TextAttributeName string = "text"

type Tag struct {
	// The name of the tag. May be the empty string
	name string
	// The inclusive starting index of the Tag - [startIndex, endIndex)
	startIdx int
	// The exclusive ending index of the Tag - [startIndex, endIndex)
	endIdx int
	// The 0-indexed depth of nesting
	depth      int
	children   []Tag
	attributes map[string]string
}

type ParseResult struct {
	// Root element containing the parsed content of the entire document
	root Tag
	// After removing leading/trailing whitespace, may not be the same slice as the input
	document []rune
}

type ParseError struct {
	startIdx int
	endIdx   int
	reason   string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("[%v,%v] %v", e.startIdx, e.endIdx, e.reason)
}

// Parse: Convert a string of tag content to a Tag tree structure (like raw HTML tags).
// It is expected that:
// - There is a single root tag
//
// Note that:
// - Leading and trailing space characters will be stripped before processing
// - Self closing tags are permitted
// - Nameless tags are permitted (but must have no attributes). i.e. </> and <>MyContent</>
// - Attributes can use either single and double quotes
// - Embedded text content will have a Tag.name of <text>.
// // i.e. For <p>Content</p>, Content will be wrapped into a tag with Tag.name = "<text>" and attribute text == "Content",
// - Empty tag attributes (valueless attributes) are not supported (i.e. <checkbox checked/>)
// - Unicode is mostly support in attribute names, values and the like
// - Escaped characters are left in-tact (i.e. &lt; won't be transformed to "<")
// - Whitespace is stripped from either side of raw text content
func Parse(runes []rune) (result ParseResult, error error) {

	if len(runes) == 0 {
		return result, &ParseError{reason: "Input in empty"}
	}

	first_none_space, last_none_space := 0, len(runes)-1
	for first_none_space < len(runes) {
		r := runes[first_none_space]
		if !unicode.IsSpace(r) {
			break
		}

		first_none_space += 1
	}

	for last_none_space > 0 {
		r := runes[last_none_space]
		if !unicode.IsSpace(r) {
			break
		}
		last_none_space -= 1
	}

	return parse(runes[first_none_space : last_none_space+1])
}

func isRuneValidForName(r rune) bool {
	return r == '_' || r == '-' || r == ':' || r == '.' || (!unicode.IsControl(r) && !unicode.IsSpace(r) && !unicode.IsPunct(r))
}

func isRuneValidForValue(r rune) bool {
	return isRuneValidForName(r) || unicode.IsPunct(r) || r == ' '
}

func parseAttributeKey(runes []rune, startIdx int) (key string, endIdx int, error error) {
	// Find Key
	currentIdx := startIdx
	for {
		r := runes[currentIdx]

		if r == '=' {
			// Found the end of the attribute key.
			key = string(runes[startIdx:currentIdx])
			return key, currentIdx, nil
		}

		if !isRuneValidForName(r) {
			// if !unicode.IsLetter(r) && (currentIdx == startIdx || (!unicode.IsNumber(r) && r != '-' && r != '_' && r != ':' && r != '.')) {
			return "", -1, &ParseError{startIdx: currentIdx, endIdx: currentIdx + 1, reason: fmt.Sprintf("Unexpected rune in attribute name - %v", string(r))}
		}

		currentIdx += 1
		if currentIdx >= len(runes) {
			return "", -1, &ParseError{startIdx: startIdx, endIdx: currentIdx, reason: "Parser reached the end of the input without completing attribute name"}
		}
	}
}

func parseAttributeValue(runes []rune, startIdx int) (value string, endIdx int, error error) {
	currentIdx := startIdx
	if runes[currentIdx] != '"' && runes[currentIdx] != '\'' {
		return "", -1, &ParseError{startIdx: currentIdx, endIdx: currentIdx + 1, reason: fmt.Sprintf("Invalid attribute value quotation. Should be \" or '. Was %v", string(runes[currentIdx]))}
	}

	var quotation = runes[currentIdx]

	currentIdx += 1
	valueStart := currentIdx
	// Find Value
	for {
		r := runes[currentIdx]

		if r == quotation {
			// Successfully parsed value. We can return the true values
			value := string(runes[valueStart:currentIdx])
			return value, currentIdx, nil
		}

		if !isRuneValidForValue(r) {
			return "", -1, &ParseError{startIdx: currentIdx, endIdx: currentIdx + 1, reason: fmt.Sprintf("Unexpected rune in attribute value - %v", string(r))}
		}

		currentIdx += 1
		if currentIdx >= len(runes) {
			return "", -1, &ParseError{startIdx: valueStart, endIdx: currentIdx, reason: "Parser reached the end of the input without finding attribute value"}
		}
	}
}

func parseAttribute(runes []rune, startIdx int) (key string, value string, endIdx int, error error) {
	currentIdx := startIdx
	key, currentIdx, error = parseAttributeKey(runes, currentIdx)
	if error != nil {
		return key, value, -1, error
	}

	// Step over the "=" rune
	currentIdx += 1

	value, currentIdx, error = parseAttributeValue(runes, currentIdx)
	if error != nil {
		return key, value, -1, error
	}

	return key, value, currentIdx, nil
}

// Parse a tag
// The first character must be a '<', but any amount of spaces are permitted in the tag
// The new tag will be added to the tag stack, and the to children of it's parent (provided this entity exists)
func parseOpeningTag(runes []rune, startIdx int, parent *Tag, depth int) (tag *Tag, exitIdx int, error error) {
	if runes[startIdx] != rune('<') {
		// Not a legitimate starting tag
		return nil, -1, &ParseError{startIdx: startIdx, endIdx: startIdx + 1, reason: fmt.Sprintf("Expected an opening tag - got %v", string(runes[startIdx]))}
	}

	if parent != nil {
		parent.children = append(parent.children, Tag{startIdx: startIdx, depth: depth})
		tag = &parent.children[len(parent.children)-1]
	} else {
		tag = &Tag{startIdx: startIdx, depth: depth}
	}

	currentIdx := startIdx + 1

	// Extract the name part of the tag
	for currentIdx < len(runes) {
		r := runes[currentIdx]
		if r == ' ' || r == '/' || r == '>' {
			// Successfully found name bounds - Exit loop
			tag.name = string(runes[startIdx+1 : currentIdx])
			break
		}

		if !isRuneValidForName(r) {
			return nil, -1, &ParseError{startIdx: currentIdx, endIdx: currentIdx + 1, reason: fmt.Sprintf("Invalid rune in tag name input: %v", r)}
		}

		currentIdx += 1
	}

	// Parse attributes and search for ending tag
	for currentIdx < len(runes) {
		r := runes[currentIdx]

		if r == ' ' {
			currentIdx += 1
			continue
		}

		if r == '/' {
			// We potentially have a self closing tag.
			switch {
			case currentIdx+1 >= len(runes):
				return nil, -1, &ParseError{startIdx: currentIdx, endIdx: currentIdx, reason: "Expected a closing tag - got end of input"}
			case runes[currentIdx+1] != '>':
				return nil, -1, &ParseError{startIdx: currentIdx, endIdx: currentIdx + 1, reason: fmt.Sprintf("Expected a closing tag - got %v", runes[currentIdx+1])}
			default:
				// Must have a valid self-closing tag. Add tag to tag stack, and then return
				endIdx := currentIdx + 2
				tag.endIdx = endIdx

				return tag, endIdx, nil
			}
		}

		if r == '>' {
			// Tag is closing - We've already parsed as much information as possible. Return.
			return tag, currentIdx + 1, nil
		}

		if !unicode.IsControl(r) {
			if tag.name == "" {
				return nil, -1, &ParseError{startIdx: startIdx, endIdx: currentIdx + 1, reason: "Nameless tags cannot contain attributes"}
			}

			// Must be adding a new attribute
			var key, value string
			key, value, currentIdx, error = parseAttribute(runes, currentIdx)
			if error != nil {
				return nil, -1, error
			}

			if tag.attributes == nil {
				tag.attributes = map[string]string{}
			}
			tag.attributes[key] = value

			// Step over closing quote
			currentIdx += 1
			if currentIdx < len(runes) {
				// If there are multiple attributes, there must be a space between them.
				// Otherwise, we need to immediately close the tag
				next_rune := runes[currentIdx]
				if next_rune == '>' || next_rune == ' ' || next_rune == '/' {
					continue
				} else {
					return nil, -1, &ParseError{startIdx: currentIdx, endIdx: currentIdx + 1, reason: "Attributes must be separated by a space"}
				}
			}
		}

		currentIdx += 1
	}

	return nil, -1, &ParseError{startIdx: currentIdx, endIdx: len(runes), reason: "Parser reached the end of the input without finding a closing angle bracket >"}
}

func parseClosingTag(runes []rune, startIdx int, opening *Tag) (exitIdx int, error error) {
	if runes[startIdx] != rune('<') {
		// Not a legitimate starting tag
		return -1, &ParseError{startIdx: startIdx, endIdx: startIdx + 1, reason: fmt.Sprintf("Expected an opening angle bracket - got %v", string(runes[startIdx]))}
	}

	if startIdx+1 >= len(runes) || runes[startIdx+1] != '/' {
		return -1, &ParseError{startIdx: startIdx + 1, endIdx: startIdx + 1, reason: "No / at start of closing tag."}
	}

	currentIdx := startIdx + 2
	// Extra the name part of the tag
	for currentIdx < len(runes) {
		r := runes[currentIdx]
		if r == ' ' || r == '>' {
			// Successfully found name bounds - Check to see if it matches opening tag
			name := string(runes[startIdx+2 : currentIdx])
			if name != opening.name {
				return -1, &ParseError{startIdx: startIdx, endIdx: currentIdx,
					reason: fmt.Sprintf("Expected a closing tag. Got %v but needed %v", name, opening.name)}
			}

			break
		}

		if !isRuneValidForName(r) {
			return -1, &ParseError{startIdx: currentIdx, endIdx: currentIdx + 1,
				reason: fmt.Sprintf("Invalid rune in tag name input: %v", r)}
		}

		currentIdx += 1
	}

	for currentIdx < len(runes) {
		r := runes[currentIdx]
		if r == ' ' {
			currentIdx += 1
			continue
		}

		if r == '>' {
			// Successfully closed tag
			opening.endIdx = currentIdx + 1
			return currentIdx + 1, nil
		}

		return -1, &ParseError{startIdx: currentIdx, endIdx: currentIdx + 1,
			reason: fmt.Sprintf("Invalid rune in closing tag %v. Expecting closing angle bracket", string(r))}
	}

	return -1, &ParseError{startIdx: currentIdx, endIdx: len(runes),
		reason: "Parser reached the end of the input without finding a closing angle bracket >"}
}

func parseRawContent(runes []rune, startIdx int) (content string, endIdx int) {
	// Raw Content matches anything which isn't a new opening bracket
	currentIdx := startIdx
	firstRealContent := -1
	lastRealContent := -1
	for currentIdx < len(runes) {
		r := runes[currentIdx]
		if r == '<' {
			break
		}

		if isRuneValidForName(r) || unicode.IsPunct(r) {
			if firstRealContent == -1 {
				firstRealContent = currentIdx
			}
			lastRealContent = currentIdx
		}

		currentIdx += 1
	}

	var runeContent []rune
	if firstRealContent != -1 {
		runeContent = runes[firstRealContent : lastRealContent+1]
	}

	return string(runeContent), currentIdx
}

func parse(runes []rune) (result ParseResult, error error) {
	if runes[0] != '<' {
		return result, &ParseError{reason: "The document must have a single root tag, and it must start from the beginning of the input"}
	}

	tagStack := make([]*Tag, 0)
	currentIdx := 0
	var rootTag *Tag
	for currentIdx < len(runes) {
		// New Tag found in input
		var previousTag *Tag
		if len(tagStack) > 0 {
			previousTag = tagStack[len(tagStack)-1]
		}

		if runes[currentIdx] == '<' {

			if !(currentIdx+1 < len(runes) && runes[currentIdx+1] == '/') {
				// Found an opening tag
				var newTag *Tag
				newTag, currentIdx, error = parseOpeningTag(runes, currentIdx, previousTag, len(tagStack))
				if error != nil {
					return
				}

				// Self closing tags don't need to be added to the tag stack
				if newTag.endIdx == 0 {
					tagStack = append(tagStack, newTag)
				}

				if rootTag == nil {
					rootTag = newTag
				}
			} else {
				// Must be a closing tag
				if previousTag == nil {
					return result, &ParseError{startIdx: currentIdx, endIdx: currentIdx, reason: "Found a closing tag with no opening tags on the tag"}
				}

				currentIdx, error = parseClosingTag(runes, currentIdx, previousTag)
				if error != nil {
					return
				}

				tagStack = tagStack[:len(tagStack)-1]

				if len(tagStack) == 0 && currentIdx < len(runes) {
					return result, &ParseError{startIdx: currentIdx, endIdx: currentIdx, reason: "Closed root tag while there was still content to parse."}
				}
			}
		} else {
			if previousTag == nil {
				return result, &ParseError{startIdx: currentIdx, endIdx: currentIdx, reason: "Attempting to parse raw content, but there is no parent tag to attach it to"}
			}

			if unicode.IsSpace(runes[currentIdx]) {
				currentIdx += 1
				continue
			}

			// Must be raw content
			var content string
			startIdx := currentIdx
			content, currentIdx = parseRawContent(runes, currentIdx)

			previousTag.children = append(previousTag.children, Tag{name: TextTagName, startIdx: startIdx, endIdx: currentIdx, depth: len(tagStack) + 1,
				attributes: map[string]string{TextAttributeName: content}})
		}
	}

	result.document = runes
	result.root = *rootTag

	return
}
