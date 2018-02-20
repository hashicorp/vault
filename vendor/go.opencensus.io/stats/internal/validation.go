package internal

const (
	MaxNameLength = 255
)

func IsPrintable(str string) bool {
	for _, r := range str {
		if !(r >= ' ' && r <= '~') {
			return false
		}
	}
	return true
}
