package swan

type defaultExtractor struct{}

func init() {
	AddExtractor(defaultExtractor{})
}

func (d defaultExtractor) Score(a *Article) int {
	return 1
}

func (d defaultExtractor) Extract(a *Article) error {
	a.Cleanup()

	return nil
}
