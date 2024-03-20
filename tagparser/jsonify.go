package tagparser

import (
	"fmt"
	"sort"
	"strings"
)

func (tag *Tag) ToJson() string {
	var sb strings.Builder

	depth := 0
	if tag == nil {
		return ""
	}

	depth = tag.Depth
	toJson(tag, &sb, depth)
	sb.WriteString("\n")
	return sb.String()
}

func toJson(tag *Tag, sb *strings.Builder, depth int) (error error) {
	if tag == nil {
		return
	}

	root_header := strings.Repeat("    ", depth)
	inner_header := root_header + "    "

	if tag.Name == "<text>" {
		text, ok := tag.Attributes["text"]
		if !ok {
			return fmt.Errorf("Text tag is missing it's required text. Cannot render.")
		}
		sb.WriteString(fmt.Sprintf("%v\"%v\"", root_header, text))
		return
	}

	sb.WriteString(fmt.Sprintf("%v{\n", root_header))

	sb.WriteString(fmt.Sprintf("%v\"_name\": \"%v\"", inner_header, tag.Name))

	if tag.Attributes != nil {
		sb.WriteString(",\n")

		// Ensure that keys are sorted for a stable output
		keys := make([]string, len(tag.Attributes))
		i := 0
		for k := range tag.Attributes {
			keys[i] = k
			i += 1
		}
		sort.Strings(keys)

		for idx, key := range keys {
			sb.WriteString(fmt.Sprintf("%v\"%v\": \"%v\"", inner_header, key, tag.Attributes[key]))
			if idx != len(tag.Attributes)-1 {
				sb.WriteString(",")
				sb.WriteString("\n")
			}
		}
	}

	if len(tag.Children) > 0 {
		sb.WriteString(",\n")

		sb.WriteString(fmt.Sprintf("%v\"_children\": [\n", inner_header))

		for idx, child := range tag.Children {
			toJson(&child, sb, depth+2)

			if idx != len(tag.Children)-1 {
				sb.WriteString(",")
			}
			sb.WriteString("\n")
		}

		sb.WriteString(fmt.Sprintf("%v]", inner_header))
	}

	sb.WriteString(fmt.Sprintf("\n%v}", root_header))

	return
}
