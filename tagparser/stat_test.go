package tagparser

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCalculateStats_WorksWithNilTag(t *testing.T) {
	var tag *Tag
	got := CalculateStats(tag)
	want := Stats{}
	if !cmp.Equal(got, want) {
		t.Errorf("CalculateStats doesn't work with nil Tag. Got %v Want %v", got, want)
	}
}

func TestCalcStats_BasicTag(t *testing.T) {
	tag := &Tag{
		Name: "Hello!",
	}

	got := CalculateStats(tag)
	want := Stats{
		TotalTags:          1,
		TagHistogram:       map[string]int{"Hello!": 1},
		AttributeHistogram: map[string]int{},
	}

	if !cmp.Equal(got, want) {
		t.Errorf("CalculateStats doesn't work with BasicTag. Got %v Want %v", got, want)
	}
}

func TestCalcStats_ComplexTags(t *testing.T) {
	tag := &Tag{
		Name:       "Cool",
		Attributes: map[string]string{"A": "B"},
		Children: []Tag{
			{
				Name:       "Cool",
				Attributes: map[string]string{"A": "C"},
				Children: []Tag{
					{
						Name:       "<text>",
						Attributes: map[string]string{"text": "🐶"},
					},
				},
			},
			{
				Name:       "Beans",
				Attributes: map[string]string{"B": "C"},
			},
		},
	}

	got := CalculateStats(tag)
	want := Stats{
		TotalTags:          3,
		TotalTextContents:  1,
		TotalAttributes:    3,
		TagHistogram:       map[string]int{"Beans": 1, "Cool": 2},
		AttributeHistogram: map[string]int{"A": 2, "B": 1},
	}

	if !cmp.Equal(got, want) {
		t.Errorf("CalculateStats doesn't work for Complex Tags. Got %v Want %v", got, want)
	}
}
