package testutils

// file di mock del generatore per interi e float
type MockRandomGenerator struct {
	NInt   int
	NFloat float64
}

// per heartbeat ed ecg
func (g *MockRandomGenerator) Intn(n int) int {
	return g.NInt
}

// temperatura, ossimetro ecc
func (g *MockRandomGenerator) Float64() float64 {
	return g.NFloat
}
