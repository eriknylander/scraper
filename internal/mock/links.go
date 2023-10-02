package mock

type LinksParserMock struct {
	ParseLinksFunc func(b []byte) ([]string, error)
}

func (lpm *LinksParserMock) ParseLinks(b []byte) ([]string, error) {
	return lpm.ParseLinksFunc(b)
}
