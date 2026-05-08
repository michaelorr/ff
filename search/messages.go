package search

type ContentMatch struct {
	Path       string
	Icon       FileIcon
	Line       string
	LineNum    int
	MatchCol   int
	MatchLen   int
	Generation int
}

type FileIcon struct {
	Icon  string
	Color string
}

type ContentBatchMsg struct {
	Matches []ContentMatch
}
