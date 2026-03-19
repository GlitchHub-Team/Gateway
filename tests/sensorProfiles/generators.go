package sensorprofilestests
//file di mock del generatore per interi e float

type MockRandomGenerator struct {
	nInt   int
	nFloat float64
}

// per heartbeat ed ecg
func (g *MockRandomGenerator) Intn(n int) int {
	return g.nInt
}

//temperatura, ossimetro ecc
func (g *MockRandomGenerator) Float64() float64 {
	return g.nFloat
}
