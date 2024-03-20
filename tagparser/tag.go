package tagparser

var TextTagName string = "<text>"
var TextAttributeName string = "text"

type Tag struct {
	// The Name of the tag. May be the empty string
	Name string
	// The inclusive starting index of the Tag - [startIndex, endIndex)
	StartIdx int
	// The exclusive ending index of the Tag - [startIndex, endIndex)
	EndIdx int
	// The 0-indexed Depth of nesting
	Depth      int
	Children   []Tag
	Attributes map[string]string
}

func (t *Tag) Render(document []rune) string {
	return string(document[t.StartIdx:t.EndIdx])
}
