package tagparser

import (
	"testing"
)

func TestToJson_WorksWithNilTag(t *testing.T) {
	var tag *Tag
	got := tag.ToJson()
	want := ""
	if got != want {
		t.Errorf("ToJson doesn't work with nil tags. Got %v Want %v", got, want)
	}
}

func TestToJson_WorksWithTextTags(t *testing.T) {
	tag := &Tag{
		Name:       "<text>",
		Attributes: map[string]string{"text": "Hello, World!"},
	}

	got := tag.ToJson()
	want := "\"Hello, World!\"\n"
	if got != want {
		t.Errorf("ToJson doesn't work with text tags. Got %v Want %v", got, want)
	}
}

func TestToJson_WorksWithSimpleTags(t *testing.T) {
	tag := &Tag{
		Name: "🐹",
	}
	got := tag.ToJson()
	want := "{\n    \"_name\": \"🐹\"\n}\n"
	if got != want {
		t.Errorf("ToJson doesn't work with simple documents. \nGot:\n%v\nWant:\n%v", got, want)
	}
}

func TestToJson_WorksWithAttributes(t *testing.T) {
	tag := &Tag{
		Name: "🐹",
		Attributes: map[string]string{
			"🦊": "🐶",
			"🥳": "😭",
		},
	}
	got := tag.ToJson()
	want := `{
    "_name": "🐹",
    "🥳": "😭",
    "🦊": "🐶"
}
`
	if got != want {
		t.Errorf("ToJson doesn't work with attributes. \nGot:\n%v\nWant:\n%v", got, want)
	}
}

func TestToJson_BasicSingleChild(t *testing.T) {
	tag := &Tag{
		Name: "Cool",
		Children: []Tag{
			{Name: "<text>", Attributes: map[string]string{"text": "Beans!"}},
		},
	}

	got := tag.ToJson()
	want := `{
    "_name": "Cool",
    "_children": [
        "Beans!"
    ]
}
`

	if got != want {
		t.Errorf("ToJson doesn't work with single children. \nGot:\n%v\nWant:\n%v", got, want)
	}
}

func TestToJson_WorksWithChildren(t *testing.T) {
	tag := &Tag{
		Name: "Cool",
		Children: []Tag{
			{Name: "<text>", Attributes: map[string]string{"text": "Beans!"}},
			{Name: "🦊", Attributes: map[string]string{"🔥": "💧"}},
		},
	}

	got := tag.ToJson()
	want := `{
    "_name": "Cool",
    "_children": [
        "Beans!",
        {
            "_name": "🦊",
            "🔥": "💧"
        }
    ]
}
`

	if got != want {
		t.Errorf("ToJson doesn't work with children. \nGot:\n%v\nWant:\n%v", got, want)
	}
}

func TestToJson_WorksWithAllFeatures(t *testing.T) {
	tag := &Tag{
		Name:       "Cool",
		Attributes: map[string]string{"Neat": "Attribute!"},
		Children: []Tag{
			{Name: "<text>", Attributes: map[string]string{"text": "Beans!"}},
			{Name: "🦊", Attributes: map[string]string{"🔥": "💧"}},
		},
	}

	got := tag.ToJson()
	want := `{
    "_name": "Cool",
    "Neat": "Attribute!",
    "_children": [
        "Beans!",
        {
            "_name": "🦊",
            "🔥": "💧"
        }
    ]
}
`

	if got != want {
		t.Errorf("ToJson doesn't work with children. \nGot:\n%v\nWant:\n%v", got, want)
	}
}
