package domain

type GatewayStatus string

const (
	Active         GatewayStatus = "active"
	Inactive       GatewayStatus = "inactive"
	Decommissioned GatewayStatus = "decommissioned"
	Stopped        GatewayStatus = "stopped"
)
