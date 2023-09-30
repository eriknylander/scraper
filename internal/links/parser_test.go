package links_test

import (
	"testing"

	"github.com/eriknylander/scraper/internal/links"
	"github.com/stretchr/testify/require"
)

const htmlWithLinks = `
<html>
	<head>
		<link href="style.css" />
		<script src="script.js"></script>
	</head>
	<body>
		<a href="page.htm">Page 1</a>
		<img src="image.jpg" />
	</body>
</html>
`

func Test_htmlParser_ParseLinks(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    []string
		wantErr bool
	}{
		{
			name:    "No links in HTML document",
			input:   []byte(`<html><head><title>Test</title></head><body><p>Test</p></body></html>`),
			want:    []string{},
			wantErr: false,
		},
		{
			name:    "Has links",
			input:   []byte(htmlWithLinks),
			want:    []string{"style.css", "script.js", "page.htm", "image.jpg"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := links.NewHTMLParser()

			got, err := p.ParseLinks(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, got)

				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
