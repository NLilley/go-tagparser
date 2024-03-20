package tagparser

import (
	"testing"
)

func TestTagRender_WorksProperly(t *testing.T) {
	input := []rune("<html><p>🐶\n🦊</p></html>")
	tag := &Tag{Name: "p", StartIdx: 6, EndIdx: 16}
	got := tag.Render(input)
	want := "<p>🐶\n🦊</p>"
	if got != want {
		t.Errorf("Tag.Render() incorrectly rendered it's content. Got '%v'\n Want '%v'", got, want)
	}
}
