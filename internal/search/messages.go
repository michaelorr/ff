package search

type ContentMatch struct {
	Path    string
	Line    string
	LineNum int
	Col     int
	Len     int
}

type ContentBatchMsg struct {
	Gen     uint64
	Matches []ContentMatch
}

type ContentDoneMsg struct {
	Gen uint64
}
