//go:build cgo

package version

func init() {
	CgoEnabled = true
}
