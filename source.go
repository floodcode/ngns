package ngns

// StringSource represents string source used to suply bruteforcer with credentials
type StringSource interface {
	// Next returns next string from a string source
	Next() (string, bool)
	// Reset resets string source
	Reset()
}

// ArraySource a StringSource that can be created using string array
type ArraySource struct {
	StringSource
	source []string
	index  int
}

// NewArraySource creates new ArraySource
func NewArraySource(source []string) *ArraySource {
	return &ArraySource{
		source: source,
		index:  0,
	}
}

// Next returns next string from array source
func (s *ArraySource) Next() (str string, ok bool) {
	if len(s.source) <= s.index {
		return "", false
	}

	res := s.source[s.index]
	s.index++
	return res, true
}

// Reset resets string source
func (s ArraySource) Reset() {
	s.index = 0
}
