// +build !windows

package manager

func defaultUploadBufferProvider() ReadSeekerWriteToProvider {
	return nil
}
