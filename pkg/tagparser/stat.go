package tagparser

import (
	"fmt"
	"strings"
)

type Stats struct {
	TotalTags          int
	TotalTextContents  int
	TotalAttributes    int
	TagHistogram       map[string]int
	AttributeHistogram map[string]int
}

func (s *Stats) Render() string {
	var builder strings.Builder

	printMap := func(m map[string]int) {
		for key, value := range m {
			builder.WriteString(fmt.Sprintf("\t%v\t%v\n", key, value))
		}
	}

	builder.WriteString("\nTag Document Statistics:\n")
	builder.WriteString("------------------------\n")
	builder.WriteString(fmt.Sprintf("Total Tags: %v\n", s.TotalTags))
	builder.WriteString(fmt.Sprintf("Total Text Contents: %v\n", s.TotalTextContents))
	builder.WriteString(fmt.Sprintf("Total Attributes: %v\n", s.TotalAttributes))

	builder.WriteString("\nTag Histogram:\n")
	printMap(s.TagHistogram)

	builder.WriteString("\nAttribute Histogram:\n")
	printMap(s.AttributeHistogram)

	return builder.String()
}

func CalculateStats(tag *Tag) (stats Stats) {
	if tag == nil {
		return
	}

	stats.TagHistogram = map[string]int{}
	stats.AttributeHistogram = map[string]int{}

	var visit func(t *Tag)
	visit = func(t *Tag) {
		stats.TotalTags += 1

		tagCount, ok := stats.TagHistogram[t.Name]
		if !ok {
			stats.TagHistogram[t.Name] = 1
		} else {
			stats.TagHistogram[t.Name] = tagCount + 1
		}

		if t.Attributes != nil {
			for key := range t.Attributes {
				stats.TotalAttributes += 1
				count, ok := stats.AttributeHistogram[key]
				if !ok {
					stats.AttributeHistogram[key] = 1
				} else {
					stats.AttributeHistogram[key] = count + 1
				}
			}
		}

		if t.Children == nil {
			return
		}

		for _, child := range t.Children {
			if child.Name == TextTagName {
				stats.TotalTextContents += 1
			} else {
				visit(&child)
			}
		}
	}

	visit(tag)

	return
}
