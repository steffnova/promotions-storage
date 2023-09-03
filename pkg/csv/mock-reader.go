package csv

type mockReader struct {
	implRead func(p []byte) (n int, err error)
}

// Read always returns an EOF error.
func (er *mockReader) Read(p []byte) (n int, err error) {
	return er.implRead(p)
}
