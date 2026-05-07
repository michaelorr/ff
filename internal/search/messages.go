package search

type ContentMatch struct {
	Path       string
	Line       string
	LineNum    int
	MatchCol   int
	MatchLen   int
	Generation int
}

type ContentBatchMsg struct {
	Matches []ContentMatch
}
