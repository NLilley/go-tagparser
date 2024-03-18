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

	test_inputs := []Def{
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

	for _, def := range test_inputs {
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

	test_inputs := []Def{
		{input: []rune("ter'=ible"), expectedError: "Unexpected rune"},
		{input: []rune("teri;=ble"), expectedError: "Unexpected rune"},
		{input: []rune("ter =ble"), expectedError: "Unexpected rune"},
		{input: []rune("te\nri=ble"), expectedError: "Unexpected rune"},
	}

	for _, def := range test_inputs {
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

	test_inputs := []Def{
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

	for _, def := range test_inputs {
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

	test_inputs := []Def{
		{input: []rune("hello='wor\nld'"), startIdx: 6, expectedError: "Unexpected rune in attribute value"},
		{input: []rune("hello=\"world'"), startIdx: 6, expectedError: "Parser reached the end"},
		{input: []rune("hello='world\""), startIdx: 6, expectedError: "Parser reached the end"},
		{input: []rune("hello=world\""), startIdx: 5, expectedError: "Invalid attribute value quotation"},
		{input: []rune("hello='world"), startIdx: 6, expectedError: "Parser reached the end"},
		{input: []rune("hello=world"), startIdx: 5, expectedError: "Invalid attribute value quotation"},
		{input: []rune("<div hello=world />"), startIdx: 10, expectedError: "Invalid attribute value quotation"},
	}

	for _, def := range test_inputs {
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

	test_inputs := []Def{
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

	for _, def := range test_inputs {
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

	test_inputs := []Def{
		{input: []rune("</>"), expectedName: ""},
		{input: []rune("< />"), expectedName: ""},
		{input: []rune("<        />"), expectedName: ""},
		{input: []rune("<hello/>"), expectedName: "hello"},
		{input: []rune("<hello />"), expectedName: "hello"},
		{input: []rune("<hello best='parser' is='best'/>"), expectedName: "hello"},
	}

	for _, def := range test_inputs {
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

	test_inputs := []Def{
		{input: []rune("</"), expectedError: "got end of input"},
		{input: []rune(" />"), expectedError: "Expected an opening tag"},
		{input: []rune("<        /"), expectedError: "got end of input"},
		{input: []rune("<hell/o>"), expectedError: "Expected a closing tag"},
		{input: []rune("h<ello />"), expectedError: "Expected an opening tag"},
		{input: []rune("<hello best='parser' is='best'/a>"), expectedError: "Expected a closing tag"},
	}

	for _, def := range test_inputs {
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

	test_inputs := []Def{
		// {input: []rune("<A B='C'>"), expectedName: "A", startIdx: 0, expectedEndIdx: 9,
		// 	expectedAttributes: map[string]string{"B": "C"}},
		// {input: []rune("<animal name='Smiles' occupation='Cat'><p></p></animal>"), expectedName: "animal", startIdx: 0, expectedEndIdx: 39,
		// 	expectedAttributes: map[string]string{"name": "Smiles", "occupation": "Cat"}},
		{input: []rune("<😊 a🤦='🦊'>"), expectedName: "😊", startIdx: 0, expectedEndIdx: 10,
			expectedAttributes: map[string]string{"a🤦": "🦊"}},
	}

	for _, def := range test_inputs {
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

		if tag.attributes == nil || len(tag.attributes) != len(def.expectedAttributes) {
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

// Things to Test
// - Spaces between attributes
// - attribute name parsing
// - attribute value parsing
// - full tags with spaces
// - Empty tags with content
