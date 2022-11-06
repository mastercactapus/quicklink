package store

type entryScanner struct {
	results []entry
	first   bool
}

type entry struct{ k, v string }

func (s *entryScanner) Next() bool {
	if s.first {
		s.results = s.results[1:]
	}

	s.first = true
	return len(s.results) > 0
}

func (s *entryScanner) Err() error { return nil }
func (s *entryScanner) Key() string {
	if len(s.results) == 0 {
		return ""
	}
	return s.results[0].k
}

func (s *entryScanner) Value() string {
	if len(s.results) == 0 {
		return ""
	}
	return s.results[0].v
}

func (s *entryScanner) Close() error { return nil }
