package models

// CPD represents the CPD current status and values
type CPD struct {
	Temp        float32
	Hum         float32
	Light       bool
	UPSStatus   string
	WarningTemp bool
	WarningUPS  bool
}

// IsWarning returns true is there is any warning, false otherwise
func (c *CPD) IsWarning() bool {
	return c.WarningTemp || c.WarningUPS
}
