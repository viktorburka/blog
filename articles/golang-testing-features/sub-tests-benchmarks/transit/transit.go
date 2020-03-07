package transit

var TransitConfig struct {
	Algorithm GraphAlgorithm
}

type GraphAlgorithm int
const (
	BreadthFirst GraphAlgorithm = iota
	DepthFirst
)

type Trasit struct {
}

func NewTransit() *Trasit {
	return &Trasit{}
}

func (t *Trasit) Direct(first, second string) bool {
	return true
}