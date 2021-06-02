// +build openbsd,arm

package diagnose

func diskUsage(ctx context.Context) error {
	SpotSkipped("disk usage", "unsupported on this platform")
}
