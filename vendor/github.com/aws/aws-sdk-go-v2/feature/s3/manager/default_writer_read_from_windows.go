package manager

func defaultDownloadBufferProvider() WriterReadFromProvider {
	return NewPooledBufferedWriterReadFromProvider(1024 * 1024)
}
