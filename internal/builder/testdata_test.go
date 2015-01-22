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

func walkdir(n *Node, path string, f func(path string, n *Node)) {
	f(path, n)
	for _, e := range n.entries {
		walkdir(e, filepath.Join(path, e.name), f)
	}
}

func makedir(t *testing.T, testdir *Node) {
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

func usingTestdir(t *testing.T, testdir *Node, f func()) {
	// set up testdir
	makedir(t, testdir)
	// clean up testdir
	defer func() {
		if err := os.RemoveAll(testdir.name); err != nil {
			t.Errorf("removedir: %v", err)
		}
	}()

	// call test function
	f()
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
											content: `{"Word": "add.", "Image": "add.jpg", "Desc": "add. desc", "Hint": "add. hint", "Comp": "PowerPC instruction" }`,
										},
										{
											name:    "branch.data",
											entries: nil,
											content: `{"Word": "branch", "Image": "branch.jpg", "Desc": "branch desc", "Hint": "branch hint", "Comp": "PowerPC instruction" }`,
										},
										{
											name:    "cards.info",
											entries: nil,
											content: `
											{
												"Display": "basic",
												"Templates": ["thisdoesx", "xdoesthis"],
												"Info": [
	                                                {"File": "add..data", "Tmpl":"thisdoesx", "Count": 7, "Bucket": 0 },
	                                                {"File": "add..data", "Tmpl":"xdoesthis", "Count": 3, "Bucket": 1 },
	                                                {"File": "branch.data", "Tmpl":"thisdoesx", "Count": 1, "Bucket": 2 },
	                                                {"File": "branch.data", "Tmpl":"xdoesthis", "Count": 0, "Bucket": 3 }
                                            	]
                                            }`,
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
											content: `{"Word": "push", "Image": "push.jpg", "Desc": "push desc", "Hint": "push hint", "Comp": "git command" }`,
										},
										{
											name:    "commit.data",
											entries: nil,
											content: `{"Word": "commit", "Image": "commit.jpg", "Desc": "commit desc", "Hint": "commit hint", "Comp": "git command" }`,
										},
										{
											name:    "cards.info",
											entries: nil,
											content: `
											{
												"Display": "basic",
												"Templates": ["thisdoesx", "xdoesthis"],
												"Info": [
	                                                {"File": "push.data", "Tmpl": "thisdoesx", "Count": 2, "Bucket": 0 },
	                                                {"File": "push.data", "Tmpl": "xdoesthis", "Count": 1, "Bucket": 1 },
	                                                {"File": "commit.data", "Tmpl": "thisdoesx", "Count": 0, "Bucket": 2 },
	                                                {"File": "commit.data", "Tmpl": "xdoesthis", "Count": 1, "Bucket": 3 }
	                                            ]
                                           	}`,
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

var nestedTestdir = &Node{
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
							name: "facts",
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
													content: `{"Word": "add.", "Image": "add.jpg", "Desc": "add. desc", "Hint": "add. hint", "Comp": "PowerPC instruction" }`,
												},
												{
													name:    "branch.data",
													entries: nil,
													content: `{"Word": "branch", "Image": "branch.jpg", "Desc": "branch desc", "Hint": "branch hint", "Comp": "PowerPC instruction" }`,
												},
												{
													name:    "cards.info",
													entries: nil,
													content: `
											{
												"Display": "basic",
												"Templates": ["thisdoesx", "xdoesthis"],
												"Info": [
	                                                {"File": "add..data", "Tmpl": "thisdoesx", "Count": 7, "Bucket": 0 },
	                                                {"File": "add..data", "Tmpl": "xdoesthis", "Count": 3, "Bucket": 1 },
	                                                {"File": "branch.data", "Tmpl": "thisdoesx", "Count": 1, "Bucket": 2 },
	                                                {"File": "branch.data", "Tmpl": "xdoesthis", "Count": 0, "Bucket": 3 }
                                            	]
                                            }`,
												},
											},
										},
										{
											name:    "image",
											entries: []*Node{},
										},
										{
											name:    "tmpl",
											entries: []*Node{},
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
