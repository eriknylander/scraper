package mock

type FileWriterMock struct {
	WriteFileFunc func(path string, data []byte) error
}

func (fwm *FileWriterMock) WriteFile(path string, data []byte) error {
	return fwm.WriteFileFunc(path, data)
}
