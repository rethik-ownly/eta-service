package geocache



type LocationDetails struct {
	CityID       string
	CityName     string
	ZoneID       string
	ZoneName     string
	LocalityID   string
	LocalityName string
}

type layerInfo struct {
	ID         string
	Name       string
	Properties *layerProperties
}

type layerProperties struct {
	CityID string
}
