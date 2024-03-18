package parser

import (
	"fmt"
	"unicode"
)

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
	root     Tag
	document string
}

type ParseError struct {
	startIdx int
	endIdx   int
	reason   string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("[%v,%v] %v", e.startIdx, e.endIdx, e.reason)
}

// Parse: Convert a string of tag content (like raw HTML tags) and converts them to a Tag tree structure.
// It is expected that:
// - There is a single root node
// - There is no padding/white space at the start or the end of the input
//
// Note that:
// - Self closing tags are permitted
// - Nameless tags are permitted (but must have no attributes). i.e. </> and <>MyContent</>
// - Attributes can use either single and double quotes
// - Embedded text content will have a Tag.name of <text>.  i.e. <p>Content</p> -> Tag.name = "text",
// - Empty tag attributes (valueless attributes) are not supported (i.e. <checkbox checked/>)
// - Unicode letters are supported in tag names, attributes names, etc. The exception being control special characters (newlines, <, {, etc.)
func Parse(input []rune) (ParseResult, error) {
	result := ParseResult{}

	if len(input) == 0 {
		return result, &ParseError{reason: "Input in empty"}
	}

	if input[0] != '<' {
		return result, &ParseError{reason: ""}
	}

	result, error := parse(input)
	if error != nil {
		return result, error
	}

	return result, nil
}

func parse(input []rune) (ParseResult, error) {
	// runes := []rune(input)
	// tagStack := make([]Tag, 0)
	// currentRune := 0

	// for currentRune < len(runes) {
	// 	return parseTag(runes, tagStack, 0)
	// }
	return ParseResult{}, nil
}

// Parse a tag
// The first character must be a '<', but any amount of spaces are permitted in the tag
// The new tag will be added to the tag stack, and the to children of it's parent (provided this entity exists)
func parseOpeningTag(runes []rune, startIdx int, parent *Tag, depth int) (tag *Tag, exitIndex int, error error) {
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
	nameEndIdx := currentIdx

	addTag := func() {
		tag.name = string(runes[startIdx+1 : nameEndIdx])
	}

	// Extra the name part of the tag
	for currentIdx < len(runes) {
		r := runes[currentIdx]
		if r == ' ' || r == '/' || r == '>' {
			// Successfully found name bounds - Exit loop
			tag.name = string(runes[startIdx+1 : currentIdx])
			nameEndIdx = currentIdx
			break
		}

		if unicode.IsControl(r) {
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

				addTag()
				return tag, endIdx, nil
			}
		}

		if r == '>' {
			// Tag is closing - We've already parsed as much information as possible. Return.
			return tag, currentIdx + 1, nil
		}

		if !unicode.IsControl(r) {
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
				}
			}
		}

		currentIdx += 1
	}

	return nil, -1, &ParseError{startIdx: currentIdx, endIdx: len(runes), reason: "Parser reached the end of the input without finding a closing angle bracket >"}
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
			return "", -1, &ParseError{startIdx: startIdx, endIdx: currentIdx, reason: "Parser reached the end of the input without finding attribute name"}
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

func isRuneValidForName(r rune) bool {
	return r == '_' || r == '-' || r == ':' || r == '.' || (!unicode.IsControl(r) && !unicode.IsSpace(r) && !unicode.IsPunct(r))
}

func isRuneValidForValue(r rune) bool {
	return isRuneValidForName(r) || unicode.IsPunct(r) || r == ' '
}
