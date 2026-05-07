package search

import (
	"bufio"
	"context"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	tea "charm.land/bubbletea/v2"
	"golang.org/x/sync/errgroup"
)

// Scanner manages the state of an ongoing search. It:
// - NewScanner spawns a goroutine to receive file paths from `walker`
// - NextCmd is invoked by `model.Init` and subsequently by `model.Update` to retrieve batches of search matches
type Scanner struct {
	mu         sync.Mutex
	paths      []string
	query      string
	generation int

	cancel   context.CancelFunc
	eg       *errgroup.Group
	egCtx    context.Context
	resultCh chan ContentMatch
}

// NewScanner initialize a new Scanner and start a goroutine to receive file paths from `walker` via `pathCh`.
func NewScanner() *Scanner {
	s := &Scanner{
		resultCh: make(chan ContentMatch, runtime.NumCPU()*2),
	}

	// Start a goroutine to receive file paths from `walker` and add them to `s.paths`.
	// If there is an active search query, also enqueue the path for scanning.
	go func() {
		for path := range pathCh {
			s.mu.Lock()
			s.paths = append(s.paths, path)
			eg, egCtx, query, gen := s.eg, s.egCtx, s.query, s.generation
			s.mu.Unlock()

			if eg != nil {
				s.queuePath(eg, egCtx, path, query, gen)
			}
		}
	}()

	return s
}

// Search initiates a new search with a given query and generation number.
// Cancel any ongoing work and start new goroutines for each file path in `s.paths` to scan for the query.
func (s *Scanner) Search(query string, generation int) {
	s.mu.Lock()
	if s.cancel != nil {
		s.cancel()
	}
	s.query = query
	s.generation = generation
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel
	eg, egCtx := errgroup.WithContext(ctx)
	eg.SetLimit(runtime.NumCPU())
	s.eg = eg
	s.egCtx = egCtx
	paths := make([]string, len(s.paths))
	copy(paths, s.paths)
	s.mu.Unlock()

	for _, path := range paths {
		s.queuePath(eg, egCtx, path, query, generation)
	}
}

// NextCmd returns a tea.Cmd that encapsulates logic for retrieving batches of matches
// from `resultCh` and batching them to be sent as a `ContentBatchMsg` to the model.
// A batch is sent when either with 40 matches or after 20 milliseconds have passed, whichever comes first.
func (s *Scanner) NextCmd() tea.Cmd {
	return func() tea.Msg {
		batch := []ContentMatch{<-s.resultCh}

		timer := time.NewTimer(20 * time.Millisecond)
		defer timer.Stop()

		for len(batch) < 40 {
			select {
			case m := <-s.resultCh:
				batch = append(batch, m)
			case <-timer.C:
				return ContentBatchMsg{Matches: batch}
			}
		}

		return ContentBatchMsg{Matches: batch}
	}
}

// Add a file path to the scanning queue.
func (s *Scanner) queuePath(eg *errgroup.Group, egCtx context.Context, path, query string, generation int) {
	eg.Go(func() error {
		if egCtx.Err() != nil {
			return egCtx.Err()
		}
		return scanFile(egCtx, path, query, generation, s.resultCh)
	})
}

const maxLineLen = 2000

// Actually scan the file for the query and emit matches to `resultCh`.
// Long lines are skipped because we dont want to return matches for minified files.
func scanFile(ctx context.Context, path, query string, generation int, resultCh chan<- ContentMatch) error {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, maxLineLen), maxLineLen)
	lineNum := 0
	for scanner.Scan() {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		lineNum++
		line := scanner.Text()
		col := strings.Index(line, query)
		if col < 0 {
			continue
		}
		resultCh <- ContentMatch{
			Path:       path,
			Line:       line,
			LineNum:    lineNum,
			MatchCol:   col,
			MatchLen:   len(query),
			Generation: generation,
		}
	}
	if scanner.Err() == bufio.ErrTooLong {
		return nil
	}
	return scanner.Err()
}
