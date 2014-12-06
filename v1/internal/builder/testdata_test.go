package builder

import (
	"os"
	"path/filepath"
	"testing"
)

/*
dir structure

html/
    static/
        css/
            style.css
        js/
            ...
        index.html
    decks/
        ppc/
            cards/
                add..data
                add.xer.data
            image/
                add.jpg
            tmpl/
                thisdoesx
                xdoeswhat
        git/
            cards/
            image/
            tmpl/
    tmpl/
        thisdoesx
        xdoeswhat
*/

type Node struct {
	name    string
	entries []*Node // nil if the entry is a file
	content string  // only files have content.
}

var testdir = &Node{
	name: "testdata",
	entries: []*Node{
		{
			name: "html",
			entries: []*Node{
				{name: "static", entries: []*Node{}},
				{name: "index.html", entries: nil},
				{
					name: "decks",
					entries: []*Node{
						{
							name: "ppc",
							entries: []*Node{
								{
									name: "cards",
									entries: []*Node{
										{
											name:    "add..data",
											entries: nil,
											content: `{"Word": "add.", "Image": "add.jpg", "Desc": "add. desc", "Hint": "add. hint", "Comp": 1 }`,
										},
										{
											name:    "branch.data",
											entries: nil,
											content: `{"Word": "branch", "Image": "branch.jpg", "Desc": "branch desc", "Hint": "branch hint", "Comp": 1 }`,
										},
										{
											name:    "cards.info",
											entries: nil,
											content: `[
                                                {"File": "add.", "Type": 1, "Count": 7, "Bucket": 0 },
                                                {"File": "add.", "Type": 2, "Count": 3, "Bucket": 1 },
                                                {"File": "branch", "Type": 1, "Count": 1, "Bucket": 2 },
                                                {"File": "branch", "Type": 2, "Count": 0, "Bucket": 3 }
                                            ]`,
										},
									},
								},
								{
									name: "image",
									entries: []*Node{
										{name: "add.jpg", entries: nil},
										{name: "branch.jpg", entries: nil},
									},
								},
							},
						},
						{
							name: "git",
							entries: []*Node{
								{
									name: "cards",
									entries: []*Node{
										{
											name:    "push.data",
											entries: nil,
											content: `{"Word": "push", "Image": "push.jpg", "Desc": "push desc", "Hint": "push hint", "Comp": 1 }`,
										},
										{
											name:    "commit.data",
											entries: nil,
											content: `{"Word": "commit", "Image": "commit.jpg", "Desc": "commit desc", "Hint": "commit hint", "Comp": 1 }`,
										},
										{
											name:    "cards.info",
											entries: nil,
											content: `[
                                                {"File": "push", "Type": 1, "Count": 2, "Bucket": 0 },
                                                {"File": "push", "Type": 2, "Count": 1, "Bucket": 1 },
                                                {"File": "commit", "Type": 1, "Count": 0, "Bucket": 2 },
                                                {"File": "commit", "Type": 2, "Count": 1, "Bucket": 3 }
                                            ]`,
										},
									},
								},
								{
									name: "image",
									entries: []*Node{
										{name: "push.jpg", entries: nil},
										{name: "commit.jpg", entries: nil},
									},
								},
								{
									name: "tmpl",
									entries: []*Node{
										{
											name:    "thisdoesx.tmpl",
											entries: nil,
											content: `{{define "front"}}local thisdoesx{{end}}{{define "back"}}{{end}}`,
										},
										{
											name:    "xdoesthis.tmpl",
											entries: nil,
											content: `{{define "front"}}local xdoesthis{{end}}{{define "back"}}{{end}}`,
										},
									},
								},
							},
						},
					},
				},
				{
					name: "tmpl",
					entries: []*Node{
						{
							name:    "thisdoesx.tmpl",
							entries: nil,
							content: `{{define "front"}}default thisdoesx{{end}}{{define "back"}}{{end}}`,
						},
						{
							name:    "xdoesthis.tmpl",
							entries: nil,
							content: `{{define "front"}}default xdoesthis{{end}}{{define "back"}}{{end}}`,
						},
					},
				},
			},
		},
	},
}

func walkdir(n *Node, path string, f func(path string, n *Node)) {
	f(path, n)
	for _, e := range n.entries {
		walkdir(e, filepath.Join(path, e.name), f)
	}
}

func makedir(t *testing.T) {
	walkdir(testdir, testdir.name, func(path string, n *Node) {
		if n.entries == nil {
			fd, err := os.Create(path)
			if err != nil {
				t.Errorf("makedir: %v", err)
				return
			}

			if n.content != "" {
				_, err := fd.Write([]byte(n.content))
				if err != nil {
					t.Errorf("write file: %v", err)
				}
			}
			fd.Close()
		} else {
			os.Mkdir(path, 0770)
		}
	})
}
