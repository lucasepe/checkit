package parser

type Item struct {
	Title string
	Notes []string
}

type Group struct {
	Title string
	Notes []string
	Items []Item
}

type CheckList struct {
	Title  string
	Groups []Group
}
