package parser

import (
	"strings"
	"testing"
)

var good_runes []rune = []rune{'a', 'b', '1', '_', '-', '.', ':', '🦊', '🎇', '🥳'}
var bad_runes []rune = []rune{'\n', '\t', ' '}

func TestIsRuneValidForName_ValidNames(t *testing.T) {
	for _, r := range good_runes {
		isValid := isRuneValidForName(r)
		if !isValid {
			t.Errorf("Expected rune to be valid, but was not!: %v", string(r))
		}
	}
}

func TestIsRuneValidForName_InvalidNames(t *testing.T) {
	for _, r := range bad_runes {
		isValid := isRuneValidForName(r)
		if isValid {
			t.Errorf("Expected rune to be invalid, but was not!: %v", string(r))
		}
	}
}

func TestIsRuneValidForValue_ValidNames(t *testing.T) {
	for _, r := range append([]rune{'.', ';', '!'}, good_runes...) {
		isValid := isRuneValidForValue(r)
		if !isValid {
			t.Errorf("Expected rune to be valid, but was not!: %v", string(r))
		}
	}
}

func TestParseAttributeKey_WorksWithValidKeys(t *testing.T) {
	type Def struct {
		input          []rune
		startIdx       int
		expectedKey    string
		expectedEndIdx int
	}

	test_defs := []Def{
		{input: []rune("hello='world'"), expectedKey: "hello", expectedEndIdx: 5},
		{input: []rune("hello=\"world\""), expectedKey: "hello", expectedEndIdx: 5},
		{input: []rune("A='B'"), expectedKey: "A", expectedEndIdx: 1},
		{input: []rune("dog='cat'"), expectedKey: "dog", expectedEndIdx: 3},
		{input: []rune("<home name='precious' />"), startIdx: 6, expectedKey: "name", expectedEndIdx: 10},
		{input: []rune("<home name=\"precious\" />"), startIdx: 6, expectedKey: "name", expectedEndIdx: 10},
		{input: []rune("<home name1=\"precious\" />"), startIdx: 6, expectedKey: "name1", expectedEndIdx: 11},
		{input: []rune("<home name_now=\"precious\" />"), startIdx: 6, expectedKey: "name_now", expectedEndIdx: 14},
		{input: []rune("<home name-now=\"precious\" />"), startIdx: 6, expectedKey: "name-now", expectedEndIdx: 14},
		{input: []rune("<home name:now=\"precious\" />"), startIdx: 6, expectedKey: "name:now", expectedEndIdx: 14},
		{input: []rune("<home name.now=\"precious\" />"), startIdx: 6, expectedKey: "name.now", expectedEndIdx: 14},
	}

	for _, def := range test_defs {
		key, end_idx, err := parseAttributeKey(def.input, def.startIdx)
		if err != nil {
			t.Errorf("Got error: %v", err)
		}
		if key != def.expectedKey {
			t.Errorf("Key was incorrect - got %v want %v", key, def.expectedKey)
		}
		if end_idx != def.expectedEndIdx {
			t.Errorf("End Index was incorrect - got %v want %v", end_idx, def.expectedEndIdx)
		}
	}
}

func TestParseAttributeKey_ErrorsWithInvalidKeys(t *testing.T) {
	type Def struct {
		input         []rune
		startIdx      int
		expectedError string
	}

	test_defs := []Def{
		{input: []rune("ter'=ible"), expectedError: "Unexpected rune"},
		{input: []rune("teri;=ble"), expectedError: "Unexpected rune"},
		{input: []rune("ter =ble"), expectedError: "Unexpected rune"},
		{input: []rune("te\nri=ble"), expectedError: "Unexpected rune"},
	}

	for _, def := range test_defs {
		_, _, err := parseAttributeKey(def.input, def.startIdx)

		if err == nil || !strings.Contains(err.Error(), def.expectedError) {
			t.Errorf("Error was %v but expected %v", err, def.expectedError)
		}
	}
}

func TestParseAttributeValue_WorksWithValidValues(t *testing.T) {
	type Def struct {
		input          []rune
		startIdx       int
		expectedValue  string
		expectedEndIdx int
	}

	test_defs := []Def{
		{input: []rune("hello='world'"), expectedValue: "world", startIdx: 6, expectedEndIdx: 12},
		{input: []rune("hello=\"world\""), expectedValue: "world", startIdx: 6, expectedEndIdx: 12},
		{input: []rune("A='B'"), expectedValue: "B", startIdx: 2, expectedEndIdx: 4},
		{input: []rune("<div A='B'>"), expectedValue: "B", startIdx: 7, expectedEndIdx: 9},
		{input: []rune("dog='cat'"), expectedValue: "cat", startIdx: 4, expectedEndIdx: 8},
		{input: []rune("dog='😊😊😊'"), expectedValue: "😊😊😊", startIdx: 4, expectedEndIdx: 8},
		{input: []rune("dog='___'"), expectedValue: "___", startIdx: 4, expectedEndIdx: 8},
		{input: []rune("dog='111'"), expectedValue: "111", startIdx: 4, expectedEndIdx: 8},
		{input: []rune("dog='11&#34;1'"), expectedValue: "11&#34;1", startIdx: 4, expectedEndIdx: 13},
	}

	for _, def := range test_defs {
		value, end_idx, err := parseAttributeValue(def.input, def.startIdx)
		if err != nil {
			t.Errorf("Got error: %v", err)
		}
		if value != def.expectedValue {
			t.Errorf("Value was incorrect - got %v want %v", value, def.expectedValue)
		}
		if end_idx != def.expectedEndIdx {
			t.Errorf("End Index was incorrect - got %v want %v", end_idx, def.expectedEndIdx)
		}
	}
}

func TestParseAttributeValue_FailsWithInvalidValues(t *testing.T) {
	type Def struct {
		input         []rune
		startIdx      int
		expectedError string
	}

	test_defs := []Def{
		{input: []rune("hello='wor\nld'"), startIdx: 6, expectedError: "Unexpected rune in attribute value"},
		{input: []rune("hello=\"world'"), startIdx: 6, expectedError: "Parser reached the end"},
		{input: []rune("hello='world\""), startIdx: 6, expectedError: "Parser reached the end"},
		{input: []rune("hello=world\""), startIdx: 5, expectedError: "Invalid attribute value quotation"},
		{input: []rune("hello='world"), startIdx: 6, expectedError: "Parser reached the end"},
		{input: []rune("hello=world"), startIdx: 5, expectedError: "Invalid attribute value quotation"},
		{input: []rune("<div hello=world />"), startIdx: 10, expectedError: "Invalid attribute value quotation"},
	}

	for _, def := range test_defs {
		_, _, err := parseAttributeValue(def.input, def.startIdx)

		if err == nil || !strings.Contains(err.Error(), def.expectedError) {
			t.Errorf("Error was %v but expected %v", err, def.expectedError)
		}
	}
}

func TestParseAttribute_WorksWithFullPair(t *testing.T) {
	type Def struct {
		input          []rune
		startIdx       int
		expectedKey    string
		expectedValue  string
		expectedEndIdx int
	}

	test_defs := []Def{
		{input: []rune("hello='world'"), expectedKey: "hello", expectedValue: "world", startIdx: 0, expectedEndIdx: 12},
		{input: []rune("hello=\"world\""), expectedKey: "hello", expectedValue: "world", startIdx: 0, expectedEndIdx: 12},
		{input: []rune("A='B'"), expectedKey: "A", expectedValue: "B", startIdx: 0, expectedEndIdx: 4},
		{input: []rune("A='B' C='D'"), expectedKey: "A", expectedValue: "B", startIdx: 0, expectedEndIdx: 4},
		{input: []rune("A='B' C='D'"), expectedKey: "C", expectedValue: "D", startIdx: 6, expectedEndIdx: 10},
		{input: []rune("<div A='B'>"), expectedKey: "A", expectedValue: "B", startIdx: 5, expectedEndIdx: 9},
		{input: []rune("dog='cat'"), expectedKey: "dog", expectedValue: "cat", startIdx: 0, expectedEndIdx: 8},
		{input: []rune("dog='😊😊😊'"), expectedKey: "dog", expectedValue: "😊😊😊", startIdx: 0, expectedEndIdx: 8},
		{input: []rune("fox🦊='likes to party 🥳'"), expectedKey: "fox🦊", expectedValue: "likes to party 🥳", startIdx: 0, expectedEndIdx: 22},
		{input: []rune("dog='111'"), expectedKey: "dog", expectedValue: "111", startIdx: 0, expectedEndIdx: 8},
		{input: []rune("dog='11&#34;1'"), expectedKey: "dog", expectedValue: "11&#34;1", startIdx: 0, expectedEndIdx: 13},
	}

	for _, def := range test_defs {
		key, value, end_idx, err := parseAttribute(def.input, def.startIdx)
		if err != nil {
			t.Errorf("Got error: %v", err)
			continue
		}
		if key != def.expectedKey {
			t.Errorf("Key was incorrect - got %v want %v", key, def.expectedKey)
		}
		if value != def.expectedValue {
			t.Errorf("Value was incorrect - got %v want %v", value, def.expectedValue)
		}
		if end_idx != def.expectedEndIdx {
			t.Errorf("End Index was incorrect - got %v want %v", end_idx, def.expectedEndIdx)
		}
	}
}

func TestParseOpeningTag_WorksWithValidSelfClosingTags(t *testing.T) {
	type Def struct {
		input        []rune
		expectedName string
	}

	test_defs := []Def{
		{input: []rune("</>"), expectedName: ""},
		{input: []rune("< />"), expectedName: ""},
		{input: []rune("<        />"), expectedName: ""},
		{input: []rune("<hello/>"), expectedName: "hello"},
		{input: []rune("<hello />"), expectedName: "hello"},
		{input: []rune("<hello best='parser' is='best'/>"), expectedName: "hello"},
	}

	for _, def := range test_defs {
		parent := &Tag{}
		tag, endIndex, err := parseOpeningTag(def.input, 0, parent, 10)
		if err != nil {
			t.Errorf("Parse failed when expected to succeed! %v", err)
		}
		if endIndex != len(def.input) {
			t.Errorf("Invalid endIndex for input %v", def)
		}
		if len(parent.children) != 0 {
			child := &parent.children[0]
			if child != tag {
				t.Errorf("Newly created tag was not added to parent")
			}
		}
		if tag.name != def.expectedName {
			t.Errorf("Newly created tag had incorrect name. Name was '%v'. Should be '%v'", tag.name, def.expectedName)
		}
		if tag.depth != 10 {
			t.Errorf("Depth not being set correctly")
		}
	}
}

func TestParseOpeningTag_BreaksWithInvalidSelfClosingTags(t *testing.T) {
	type Def struct {
		input         []rune
		expectedError string
	}

	test_defs := []Def{
		{input: []rune("</"), expectedError: "got end of input"},
		{input: []rune(" />"), expectedError: "Expected an opening tag"},
		{input: []rune("<        /"), expectedError: "got end of input"},
		{input: []rune("<hell/o>"), expectedError: "Expected a closing tag"},
		{input: []rune("h<ello />"), expectedError: "Expected an opening tag"},
		{input: []rune("<hello best='parser' is='best'/a>"), expectedError: "Expected a closing tag"},
	}

	for _, def := range test_defs {
		_, _, err := parseOpeningTag(def.input, 0, nil, 0)

		if err == nil || !strings.Contains(err.Error(), def.expectedError) {
			t.Errorf("Error was %v but expected %v", err, def.expectedError)
		}
	}
}

func TestParseOpeningTag_CorrectlyParsesNormalOpeningTags(t *testing.T) {
	type Def struct {
		input              []rune
		startIdx           int
		expectedName       string
		expectedEndIdx     int
		expectedAttributes map[string]string
	}

	test_defs := []Def{
		{input: []rune("<A B='C'>"), expectedName: "A", startIdx: 0, expectedEndIdx: 9,
			expectedAttributes: map[string]string{"B": "C"}},
		{input: []rune("<animal name='Smiles' occupation='Cat'><p></p></animal>"), expectedName: "animal", startIdx: 0, expectedEndIdx: 39,
			expectedAttributes: map[string]string{"name": "Smiles", "occupation": "Cat"}},
		{input: []rune("<😊 a🤦='🦊'>"), expectedName: "😊", startIdx: 0, expectedEndIdx: 10,
			expectedAttributes: map[string]string{"a🤦": "🦊"}},
		{input: []rune("<>< >Lonely</></ >"), expectedName: "", startIdx: 2, expectedEndIdx: 5},
	}

	for _, def := range test_defs {
		tag, endIdx, error := parseOpeningTag(def.input, def.startIdx, nil, 1)

		if error != nil {
			t.Errorf("Expected tag to parse correctly, but got error %v", error)
			continue
		}

		if tag.endIdx != 0 {
			t.Errorf("Un-closed tag has a non-0 end index set")
		}

		if tag.name != def.expectedName {
			t.Errorf("Got tag name %v expected %v", tag.name, def.expectedName)
		}

		if endIdx != def.expectedEndIdx {
			t.Errorf("Got end index %v expected %v", endIdx, def.expectedEndIdx)
		}

		if (tag.attributes == nil && def.expectedAttributes != nil) ||
			(tag.attributes != nil && def.expectedAttributes == nil) ||
			len(tag.attributes) != len(def.expectedAttributes) {
			t.Errorf("Invalid number of attributes in parsed result")
		}

		for key, value := range def.expectedAttributes {
			gotVal, ok := tag.attributes[key]
			if !ok || gotVal != value {
				t.Errorf("Got attribute value %v when expecting %v", gotVal, value)
			}
		}
	}
}

func TestParseOpeningTag_BreaksWithInvalidOpeningTags(t *testing.T) {
	type Def struct {
		input         []rune
		expectedError string
	}

	test_defs := []Def{
		{input: []rune("<"), expectedError: "Parser reached the end of the input"},
		{input: []rune(" >"), expectedError: "Expected an opening tag"},
		{input: []rune("<        "), expectedError: "Parser reached the end of the input"},
		{input: []rune("<hell o>"), expectedError: "without completing attribute name"},
		{input: []rune("h<ello >"), expectedError: "Expected an opening tag"},
		{input: []rune("<hello best='parser'is='best'>"), expectedError: "Attributes must be separated"},
		{input: []rune("< bad='attribute'></>"), expectedError: "Nameless tags cannot contain attributes"},
	}

	for _, def := range test_defs {
		_, endIdx, err := parseOpeningTag(def.input, 0, nil, 0)

		if err == nil || !strings.Contains(err.Error(), def.expectedError) {
			t.Errorf("Error was %v but expected %v", err, def.expectedError)
			continue
		}

		if endIdx != -1 {
			t.Errorf("Test errored, but end idx was not -1. Got %v", endIdx)
		}
	}
}

func TestParseClosingTag_WorksWithValidSetup(t *testing.T) {
	type Def struct {
		input          []rune
		tagName        string
		startIdx       int
		expectedEndIdx int
	}

	test_defs := []Def{
		{input: []rune("<a>Hello</a>"), tagName: "a", startIdx: 8, expectedEndIdx: 12},
		{input: []rune("<a>Hello</a     >"), tagName: "a", startIdx: 8, expectedEndIdx: 17},
		{input: []rune("<🐶>Woof!</🐶>"), tagName: "🐶", startIdx: 8, expectedEndIdx: 12},
		{input: []rune("<>< >Lonely</></ >"), tagName: "", startIdx: 11, expectedEndIdx: 14},
	}

	for _, def := range test_defs {
		openingTag := &Tag{name: def.tagName}
		endIdx, error := parseClosingTag(def.input, def.startIdx, openingTag)

		if error != nil {
			t.Errorf("Expected closing tag to parse correctly, but got error %v", error)
			continue
		}

		if endIdx != openingTag.endIdx || endIdx != def.expectedEndIdx {
			t.Errorf("Expected end indices to be %v but they were %v and %v", def.expectedEndIdx, openingTag.endIdx, endIdx)
		}
	}
}

func TestParseClosingTag_BreaksWithInvalidTags(t *testing.T) {
	type Def struct {
		input         []rune
		tagName       string
		expectedError string
	}

	test_defs := []Def{
		{input: []rune("<"), tagName: "", expectedError: "No /"},
		{input: []rune(" >"), tagName: "", expectedError: "Expected an opening angle bracket"},
		{input: []rune("<        "), tagName: "", expectedError: "No /"},
		{input: []rune("<hell o>"), tagName: "hello", expectedError: "No /"},
		{input: []rune("h<ello >"), tagName: "ello", expectedError: "Expected an opening angle bracket"},
		{input: []rune("</hello best='parser'is='best'>"), tagName: "hello", expectedError: "Invalid rune in closing tag"},
	}

	for _, def := range test_defs {
		openingTag := &Tag{name: def.tagName}
		endIdx, err := parseClosingTag(def.input, 0, openingTag)

		if err == nil || !strings.Contains(err.Error(), def.expectedError) {
			t.Errorf("Error was %v but expected %v", err, def.expectedError)
		}

		if endIdx != -1 {
			t.Errorf("Expecting endIdx to be -1 as input should error. Got %v", endIdx)
		}
	}
}

func TestParseRawContent_SuccessfullyDrainsContent(t *testing.T) {
	type Def struct {
		input           []rune
		startIdx        int
		expectedContent string
		expectedEndIdx  int
	}

	test_defs := []Def{
		{input: []rune(""), expectedContent: "", expectedEndIdx: 0},
		{input: []rune("Woof! 🐶"), startIdx: 0, expectedContent: "Woof! 🐶", expectedEndIdx: 7},
		{input: []rune("<p>Super!</p>"), startIdx: 3, expectedContent: "Super!", expectedEndIdx: 9},
		{input: []rune("<p>    Super!\n\n\n\n</p>"), startIdx: 3, expectedContent: "Super!", expectedEndIdx: 17},
	}

	for _, def := range test_defs {
		content, endIdx := parseRawContent(def.input, def.startIdx)
		if endIdx != def.expectedEndIdx {
			t.Errorf("EndIdx doesn't match: Got %v, want %v", endIdx, def.expectedEndIdx)
		}

		if content != def.expectedContent {
			t.Errorf("Parsed content does not match! got %v want %v", content, def.expectedContent)
		}
	}
}

func TestParse_WorksWithValidDocument_Simple(t *testing.T) {
	input := []rune("<p>Hello, World! 🐶</p>")
	result, error := Parse(input)

	if error != nil {
		t.Errorf("Expected Parse to succeeded, but it failed with error %v", error)
		return
	}

	root := &result.root
	if root.name != "p" {
		t.Errorf("Root tag had name %v but had name %v", root.name, "p")
	}

	if len(root.children) != 1 {
		t.Errorf("Expected root tag to have a single child. It had %v children", len(root.children))
	}

	child := &result.root.children[0]
	if child.name != "<text>" {
		t.Errorf("Expected single child to be a pseudo <text> tag. Got %v", child.name)
	}

	content, ok := child.attributes["text"]
	if !ok || content != "Hello, World! 🐶" {
		t.Errorf("Single child has invalid content. Got %v", content)
	}
}

func TestParse_WorksWithEmptyTags(t *testing.T) {
	input := []rune("<><>Lonely!</></>")
	result, error := Parse(input)

	if error != nil {
		t.Errorf("Expected Parse to succeeded, but it failed with error %v", error)
		return
	}

	root := &result.root
	if root.name != "" {
		t.Errorf("Root tag had name %v but had name %v", root.name, "p")
	}

	if len(root.children) != 1 {
		t.Errorf("Expected root tag to have a single child. It had %v children", len(root.children))
		return
	}

	child := &result.root.children[0]
	if child.name != "" {
		t.Errorf("Expected single child to be a nameless tag. Got %v", child.name)
	}

	if len(child.children) != 1 {
		t.Errorf("Expected root to have a single grandchild")
		return
	}

	grandChild := &child.children[0]
	if grandChild.name != "<text>" {
		t.Errorf("Expected single grandchild to a raw text tag. Got %v", grandChild.name)
	}

	content, ok := grandChild.attributes["text"]
	if !ok || content != "Lonely!" {
		t.Errorf("Payload Incorrect! Got '%v' but want '%v'", content, "Lonely!")
	}
}

func TestParse_WorksWithCompoundObject(t *testing.T) {
	var builder strings.Builder
	builder.WriteString("       <html lang='en'>\n")
	builder.WriteString("<head />\n")
	builder.WriteString("<body>\n")
	builder.WriteString("<p style='font-weight:bold;'>Cool</p>")
	builder.WriteString("<p>Beans!</p>")
	builder.WriteString("</body>\n")
	builder.WriteString("</html>      ")

	result, error := Parse([]rune(builder.String()))
	if error != nil {
		t.Errorf("Expected Parse to succeed. Got error: %v", error)
		return
	}

	root := &result.root
	if root.name != "html" {
		t.Errorf("Expected root to be a HTML element. Got %v", root.name)
	}

	html_lang, ok := root.attributes["lang"]
	if !ok || html_lang != "en" {
		t.Errorf("Expected root to have attribute lang='en'. Got %v", html_lang)
	}

	if len(root.children) != 2 {
		t.Errorf("Root should have 2 child elements. Got %v", len(root.children))
		return
	}

	head := root.children[0]
	body := root.children[1]

	if head.name != "head" {
		t.Errorf("Head element has wrong name. Got %v", head.name)
	}

	if len(head.children) > 0 {
		t.Errorf("Expecting head to have no children. Got %v", head.children)
	}

	if body.name != "body" {
		t.Errorf("Body element has wrong name. Got %v", body.name)
	}

	if len(body.children) != 2 {
		t.Errorf("Body element should have 2 children. It has %v", len(body.children))
		return
	}

	first_p := body.children[0]
	second_p := body.children[1]

	if first_p.name != "p" || second_p.name != "p" {
		t.Errorf("p tags parsed with incorrect names. Got %v and %v", first_p.name, second_p.name)
	}

	first_p_style, ok := first_p.attributes["style"]
	if !ok || first_p_style != "font-weight:bold;" {
		t.Errorf("first p style attribute parsed incorrectly")
	}

	if second_p.attributes != nil {
		t.Errorf("second p has attributes when it should not")
	}

	first_content := first_p.children[0].attributes["text"]
	if first_content != "Cool" {
		t.Errorf("first p content is incorrect. Got %v want %v", first_content, "Cool")
	}

	second_content := second_p.children[0].attributes["text"]
	if second_content != "Beans!" {
		t.Errorf("second p content is correct. Got %v want %v", second_content, "Beans!")
	}
}

func TestParse_IgnoresSurroundingSpace(t *testing.T) {
	input := []rune("\t\t\t\t\t     <>Hello!</>\n\n\n\r\n")
	result, error := Parse(input)

	if error != nil {
		t.Errorf("Expecting parse to succeed but it failed with error: %v", error)
	}

	content := result.root.children[0].attributes["text"]
	if content != "Hello!" {
		t.Errorf("Content parsed incorrectly. Got %v but want %v", content, "Hello!")
	}
}
