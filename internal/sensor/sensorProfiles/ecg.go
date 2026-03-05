package sensorprofiles

type EcgProfile struct{}

func NewEcgProfile() *EcgProfile {
	return &EcgProfile{}
}

func (g *EcgProfile) Generate() []byte {
	// Logic to generate ECG data
	return []byte("ECG data")
}
