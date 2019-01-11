package object

import "fmt"

type Device struct {
	Id               string
	IsActive         bool `json:"is_active"`
	IsPrivateSession bool `json:"is_private_session"`
	IsRestricted     bool `json:"is_restricted"`
	Name             string
	Type             string
	VolumePercent    int `json:"volume_percent"`
}

func (device *Device) String() string {
	return fmt.Sprintf(`{
  Id: %s
  IsActive: %t
  IsPrivateSession: %t
  IsRestricted: %t
  Name: %s
  Type: %s
  VolumePercent: %d
}`,
		device.Id,
		device.IsActive,
		device.IsPrivateSession,
		device.IsRestricted,
		device.Name,
		device.Type,
		device.VolumePercent,
	)
}

type Devices struct {
	Devices []Device
}
