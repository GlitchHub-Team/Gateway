package sensorprofiles

type HeartRateProfile struct{}

func NewHeartRateProfile() *HeartRateProfile {
	return &HeartRateProfile{}
}

func (g *HeartRateProfile) Generate() []byte {
	// Logic to generate heart rate data
	return []byte("heart rate data")
}
