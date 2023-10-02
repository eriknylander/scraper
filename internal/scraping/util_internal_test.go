package scraping

import "testing"

func Test_normalizePath(t *testing.T) {
	type args struct {
		url  string
		link string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "URL is root",
			args: args{
				url:  "",
				link: "static/oscar/favicon.ico",
			},
			want: "static/oscar/favicon.ico",
		},
		{
			name: "URL is subpage",
			args: args{
				url:  "catalogue/page-2.html",
				link: "category/books_1/index.html",
			},
			want: "catalogue/category/books_1/index.html",
		},
		{
			name: "URL is deeply nested, link refers to parent folder",
			args: args{
				url:  "catalogue/category/books/add-a-comment_18/index.html",
				link: "../../../../static/oscar/favicon.ico",
			},
			want: "static/oscar/favicon.ico",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizePath(tt.args.url, tt.args.link); got != tt.want {
				t.Errorf("normalizePath() = %v, want %v", got, tt.want)
			}
		})
	}
}
