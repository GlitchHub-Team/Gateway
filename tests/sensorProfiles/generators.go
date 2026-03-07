package sensorprofilestests

type MockRandomGenerator struct {
	nInt   int
	nFloat float64
}

func (g *MockRandomGenerator) Intn(n int) int {
	return g.nInt
}

func (g *MockRandomGenerator) Float64() float64 {
	return g.nFloat
}
