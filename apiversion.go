package oneview

// Supported REST API version range for HPE OneView appliances.
const (
	MinApplianceAPI = 3800
	MaxApplianceAPI = 8800

	// GlobalDashboardAPI is the current version documented in swagger 300.json.
	GlobalDashboardAPI    = 300
	MinGlobalDashboardAPI = 2
)

// Named OneView REST API versions (X-API-Version) mapped to product releases.
const (
	API3800 = 3800 // OneView 6.60
	API4000 = 4000 // OneView 7.00
	API4200 = 4200 // OneView 7.10
	API4400 = 4400 // OneView 7.20
	API4600 = 4600 // OneView 8.00
	API4800 = 4800 // OneView 8.10
	API5000 = 5000 // OneView 8.20
	API5200 = 5200 // OneView 8.30
	API5400 = 5400 // OneView 8.40
	API5600 = 5600 // OneView 8.50
	API5800 = 5800 // OneView 8.60
	API6000 = 6000 // OneView 8.70
	API6200 = 6200 // OneView 8.80
	API6400 = 6400 // OneView 8.90
	API6600 = 6600 // OneView 9.00
	API6800 = 6800 // OneView 9.10
	API7000 = 7000 // OneView 9.20
	API7200 = 7200 // OneView 9.30
	API7400 = 7400 // OneView 9.40
	API7600 = 7600 // OneView 10.00
	API7800 = 7800 // OneView 10.10
	API8000 = 8000 // OneView 10.20
	API8200 = 8200 // OneView 11.00 / 11.10 / 11.20
	API8400 = 8400
	API8600 = 8600 // OneView 11.30
	API8800 = 8800
)

// Product identifies which appliance family the client is talking to.
type Product int

const (
	ProductUnknown Product = iota
	ProductAppliance
	ProductGlobalDashboard
)

func (p Product) String() string {
	switch p {
	case ProductAppliance:
		return "oneview-appliance"
	case ProductGlobalDashboard:
		return "global-dashboard"
	default:
		return "unknown"
	}
}

// ProductRelease is a OneView product version and its REST API number.
type ProductRelease struct {
	Product string
	API     int
}

// ApplianceReleases lists known OneView appliance releases in the 3800–8800 range.
var ApplianceReleases = []ProductRelease{
	{Product: "6.60", API: API3800},
	{Product: "7.00", API: API4000},
	{Product: "7.10", API: API4200},
	{Product: "7.20", API: API4400},
	{Product: "8.00", API: API4600},
	{Product: "8.10", API: API4800},
	{Product: "8.20", API: API5000},
	{Product: "8.30", API: API5200},
	{Product: "8.40", API: API5400},
	{Product: "8.50", API: API5600},
	{Product: "8.60", API: API5800},
	{Product: "8.70", API: API6000},
	{Product: "8.80", API: API6200},
	{Product: "8.90", API: API6400},
	{Product: "9.00", API: API6600},
	{Product: "9.10", API: API6800},
	{Product: "9.20", API: API7000},
	{Product: "9.30", API: API7200},
	{Product: "9.40", API: API7400},
	{Product: "10.00", API: API7600},
	{Product: "10.10", API: API7800},
	{Product: "10.20", API: API8000},
	{Product: "11.00", API: API8200},
	{Product: "11.10", API: API8200},
	{Product: "11.20", API: API8200},
	{Product: "11.30", API: API8600},
}

func detectProduct(current int) Product {
	if current >= MinApplianceAPI {
		return ProductAppliance
	}
	if current >= MinGlobalDashboardAPI && current <= GlobalDashboardAPI {
		return ProductGlobalDashboard
	}
	return ProductUnknown
}

func isSupportedAPI(version int) bool {
	if version >= MinApplianceAPI && version <= MaxApplianceAPI {
		return true
	}
	if version >= MinGlobalDashboardAPI && version <= GlobalDashboardAPI {
		return true
	}
	return false
}
