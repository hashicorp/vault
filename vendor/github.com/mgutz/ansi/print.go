package ansi

// PrintStyles prints all style combinations to the terminal.
func PrintStyles() {
	oldPlain := plain
	plain = false

	bgColors := []string{
		"",
		":black",
		":red",
		":green",
		":yellow",
		":blue",
		":magenta",
		":cyan",
		":white",
	}
	for fg := range Colors {
		for _, bg := range bgColors {
			println(padColor(fg, []string{"" + bg, "+b" + bg, "+bh" + bg, "+u" + bg}))
			println(padColor(fg, []string{"+uh" + bg, "+B" + bg, "+Bb" + bg /* backgrounds */, "" + bg + "+h"}))
			println(padColor(fg, []string{"+b" + bg + "+h", "+bh" + bg + "+h", "+u" + bg + "+h", "+uh" + bg + "+h"}))
		}
	}
	plain = oldPlain
}

func pad(s string, length int) string {
	for len(s) < length {
		s += " "
	}
	return s
}

func padColor(s string, styles []string) string {
	buffer := ""
	for _, style := range styles {
		buffer += Color(pad(s+style, 20), s+style)
	}
	return buffer
}
