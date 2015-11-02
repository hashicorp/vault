package columnize

import "testing"

func TestListOfStringsInput(t *testing.T) {
	input := []string{
		"Column A | Column B | Column C",
		"x | y | z",
	}

	config := DefaultConfig()
	output := Format(input, config)

	expected := "Column A  Column B  Column C\n"
	expected += "x         y         z"

	if output != expected {
		t.Fatalf("\nexpected:\n%s\n\ngot:\n%s", expected, output)
	}
}

func TestEmptyLinesOutput(t *testing.T) {
	input := []string{
		"Column A | Column B | Column C",
		"",
		"x | y | z",
	}

	config := DefaultConfig()
	output := Format(input, config)

	expected := "Column A  Column B  Column C\n"
	expected += "\n"
	expected += "x         y         z"

	if output != expected {
		t.Fatalf("\nexpected:\n%s\n\ngot:\n%s", expected, output)
	}
}

func TestLeadingSpacePreserved(t *testing.T) {
	input := []string{
		"| Column B | Column C",
		"x | y | z",
	}

	config := DefaultConfig()
	output := Format(input, config)

	expected := "   Column B  Column C\n"
	expected += "x  y         z"

	if output != expected {
		t.Fatalf("\nexpected:\n%s\n\ngot:\n%s", expected, output)
	}
}

func TestColumnWidthCalculator(t *testing.T) {
	input := []string{
		"Column A | Column B | Column C",
		"Longer than A | Longer than B | Longer than C",
		"short | short | short",
	}

	config := DefaultConfig()
	output := Format(input, config)

	expected := "Column A       Column B       Column C\n"
	expected += "Longer than A  Longer than B  Longer than C\n"
	expected += "short          short          short"

	if output != expected {
		t.Fatalf("\nexpected:\n%s\n\ngot:\n%s", expected, output)
	}
}

func TestVariedInputSpacing(t *testing.T) {
	input := []string{
		"Column A       |Column B|    Column C",
		"x|y|          z",
	}

	config := DefaultConfig()
	output := Format(input, config)

	expected := "Column A  Column B  Column C\n"
	expected += "x         y         z"

	if output != expected {
		t.Fatalf("\nexpected:\n%s\n\ngot:\n%s", expected, output)
	}
}

func TestUnmatchedColumnCounts(t *testing.T) {
	input := []string{
		"Column A | Column B | Column C",
		"Value A | Value B",
		"Value A | Value B | Value C | Value D",
	}

	config := DefaultConfig()
	output := Format(input, config)

	expected := "Column A  Column B  Column C\n"
	expected += "Value A   Value B\n"
	expected += "Value A   Value B   Value C   Value D"

	if output != expected {
		t.Fatalf("\nexpected:\n%s\n\ngot:\n%s", expected, output)
	}
}

func TestAlternateDelimiter(t *testing.T) {
	input := []string{
		"Column | A % Column | B % Column | C",
		"Value A % Value B % Value C",
	}

	config := DefaultConfig()
	config.Delim = "%"
	output := Format(input, config)

	expected := "Column | A  Column | B  Column | C\n"
	expected += "Value A     Value B     Value C"

	if output != expected {
		t.Fatalf("\nexpected:\n%s\n\ngot:\n%s", expected, output)
	}
}

func TestAlternateSpacingString(t *testing.T) {
	input := []string{
		"Column A | Column B | Column C",
		"x | y | z",
	}

	config := DefaultConfig()
	config.Glue = "    "
	output := Format(input, config)

	expected := "Column A    Column B    Column C\n"
	expected += "x           y           z"

	if output != expected {
		t.Fatalf("\nexpected:\n%s\n\ngot:\n%s", expected, output)
	}
}

func TestSimpleFormat(t *testing.T) {
	input := []string{
		"Column A | Column B | Column C",
		"x | y | z",
	}

	output := SimpleFormat(input)

	expected := "Column A  Column B  Column C\n"
	expected += "x         y         z"

	if output != expected {
		t.Fatalf("\nexpected:\n%s\n\ngot:\n%s", expected, output)
	}
}

func TestAlternatePrefixString(t *testing.T) {
	input := []string{
		"Column A | Column B | Column C",
		"x | y | z",
	}

	config := DefaultConfig()
	config.Prefix = "  "
	output := Format(input, config)

	expected := "  Column A  Column B  Column C\n"
	expected += "  x         y         z"

	if output != expected {
		t.Fatalf("\nexpected:\n%s\n\ngot:\n%s", expected, output)
	}
}

func TestEmptyFieldReplacement(t *testing.T) {
	input := []string{
		"Column A | Column B | Column C",
		"x | | z",
	}

	config := DefaultConfig()
	config.Empty = "<none>"
	output := Format(input, config)

	expected := "Column A  Column B  Column C\n"
	expected += "x         <none>    z"

	if output != expected {
		t.Fatalf("\nexpected:\n%s\n\ngot:\n%s", expected, output)
	}
}

func TestEmptyConfigValues(t *testing.T) {
	input := []string{
		"Column A | Column B | Column C",
		"x | y | z",
	}

	config := Config{}
	output := Format(input, &config)

	expected := "Column A  Column B  Column C\n"
	expected += "x         y         z"

	if output != expected {
		t.Fatalf("\nexpected:\n%s\n\ngot:\n%s", expected, output)
	}
}

func TestMergeConfig(t *testing.T) {
	conf1 := &Config{Delim: "a", Glue: "a", Prefix: "a", Empty: "a"}
	conf2 := &Config{Delim: "b", Glue: "b", Prefix: "b", Empty: "b"}
	conf3 := &Config{Delim: "c", Prefix: "c"}

	m := MergeConfig(conf1, conf2)
	if m.Delim != "b" || m.Glue != "b" || m.Prefix != "b" || m.Empty != "b" {
		t.Fatalf("bad: %#v", m)
	}

	m = MergeConfig(conf1, conf3)
	if m.Delim != "c" || m.Glue != "a" || m.Prefix != "c" || m.Empty != "a" {
		t.Fatalf("bad: %#v", m)
	}

	m = MergeConfig(conf1, nil)
	if m.Delim != "a" || m.Glue != "a" || m.Prefix != "a" || m.Empty != "a" {
		t.Fatalf("bad: %#v", m)
	}

	m = MergeConfig(conf1, &Config{})
	if m.Delim != "a" || m.Glue != "a" || m.Prefix != "a" || m.Empty != "a" {
		t.Fatalf("bad: %#v", m)
	}
}
