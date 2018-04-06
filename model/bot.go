package model

type UserLocationInfo struct {
	Location       Coordinates
	LastCommand    string
	LastSearchTerm string
}

func (info UserLocationInfo) IsEmpty() bool {
	return info.Location.Latitude == 0 &&
		info.Location.Longitude == 0 &&
		info.LastCommand == "" &&
		info.LastSearchTerm == ""
}
