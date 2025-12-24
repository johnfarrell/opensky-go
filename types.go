package opensky

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
	"time"
)

// StateResponse represents the top-level JSON response from the /states/all endpoint.
type StateResponse struct {
	// Time is the time which the state vectors in this response are associated with.
	// All vectors represent the state of a vehicle at this time. It is a Unix timestamp (seconds).
	Time int64 `json:"time"`
	// States is a list of state vectors for all visible aircraft.
	States []StateVector `json:"states"`
}

// PrettyPrintStateResponse formats and prints the StateResponse to the given writer (e.g., os.Stdout)
// in a human-readable table format.
func PrettyPrintStateResponse(out io.Writer, response *StateResponse) {
	if response == nil {
		fmt.Fprintln(out, "No response data to display.")
		return
	}

	// Print the header information
	reportTime := time.Unix(response.Time, 0).UTC()
	fmt.Fprintf(out, "--- OpenSky State Report ---\n")
	fmt.Fprintf(out, "Report Time: %s\n", reportTime.Format("2006-01-02 15:04:05 UTC"))
	fmt.Fprintf(out, "Aircraft Found: %d\n\n", len(response.States))

	if len(response.States) == 0 {
		return
	}

	// Initialize a tab-writer for aligned columns.
	// The parameters are: output, minwidth, tabwidth, padding, padchar, flags
	w := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	defer w.Flush()

	// Print table header
	header := []string{
		"ICAO24", "CALLSIGN", "COUNTRY", "ON GROUND", "LAT", "LON", "ALT (m)", "VEL (m/s)", "V/R (m/s)", "HEADING", "SQUAWK", "CATEGORY",
	}
	fmt.Fprintln(w, strings.Join(header, "\t"))
	fmt.Fprintln(w, strings.Repeat("-\t", len(header))) // Separator line

	// Print each state vector as a row in the table
	for _, state := range response.States {
		// Prepare each field, handling nil pointers gracefully
		callsign := "N/A"
		if state.Callsign != nil {
			callsign = strings.TrimSpace(*state.Callsign)
		}

		latitude := "N/A"
		if state.Latitude != nil {
			latitude = fmt.Sprintf("%.4f", *state.Latitude)
		}

		longitude := "N/A"
		if state.Longitude != nil {
			longitude = fmt.Sprintf("%.4f", *state.Longitude)
		}

		baroAltitude := "N/A"
		if state.BaroAltitude != nil {
			baroAltitude = fmt.Sprintf("%.2f", *state.BaroAltitude)
		}

		velocity := "N/A"
		if state.Velocity != nil {
			velocity = fmt.Sprintf("%.2f", *state.Velocity)
		}

		verticalRate := "N/A"
		if state.VerticalRate != nil {
			verticalRate = fmt.Sprintf("%.2f", *state.VerticalRate)
		}

		heading := "N/A"
		if state.TrueTrack != nil {
			heading = fmt.Sprintf("%.1f°", *state.TrueTrack)
		}

		squawk := "N/A"
		if state.Squawk != nil {
			squawk = *state.Squawk
		}

		row := []string{
			state.Icao24,
			callsign,
			state.OriginCountry,
			fmt.Sprintf("%v", state.OnGround),
			latitude,
			longitude,
			baroAltitude,
			velocity,
			verticalRate,
			heading,
			squawk,
			state.Category.String(),
		}

		fmt.Fprintln(w, strings.Join(row, "\t"))
	}
}

type BoundingBox struct {
	LongitudeMin float64
	LongitudeMax float64
	LatitudeMin  float64
	LatitudeMax  float64
}

// PositionSource is an enumeration for the origin of the state's position information.
type PositionSource int

const (
	PositionSourceADS_B   PositionSource = 0 // from ADS-B
	PositionSourceASTERIX PositionSource = 1 // from ASTERIX
	PositionSourceMLAT    PositionSource = 2 // from MLAT
	PositionSourceFLARM   PositionSource = 3 // from FLARM
)

func (ps PositionSource) String() string {
	switch ps {
	case PositionSourceADS_B:
		return "ADS-B"
	case PositionSourceASTERIX:
		return "ASTERIX"
	case PositionSourceMLAT:
		return "MLAT"
	case PositionSourceFLARM:
		return "FLARM"
	default:
		return "Unknown"
	}
}

// AircraftCategory is an enumeration for the type of vehicle.
type AircraftCategory int

const (
	CategoryNoInformation           AircraftCategory = 0
	CategoryNoADSB                  AircraftCategory = 1
	CategoryLight                   AircraftCategory = 2 // < 15500 lbs
	CategorySmall                   AircraftCategory = 3 // 15500 to 75000 lbs
	CategoryLarge                   AircraftCategory = 4 // 75000 to 300000 lbs
	CategoryHighVortexLarge         AircraftCategory = 5
	CategoryHeavy                   AircraftCategory = 6 // > 300000 lbs
	CategoryHighPerformance         AircraftCategory = 7
	CategoryRotorcraft              AircraftCategory = 8
	CategoryGlider                  AircraftCategory = 9
	CategoryLighterThanAir          AircraftCategory = 10
	CategoryParachutist             AircraftCategory = 11
	CategoryUltraLight              AircraftCategory = 12
	CategoryReserved                AircraftCategory = 13
	CategoryUnmannedAerialVehicle   AircraftCategory = 14
	CategorySpaceVehicle            AircraftCategory = 15
	CategorySurfaceVehicleEmergency AircraftCategory = 16
	CategorySurfaceVehicleService   AircraftCategory = 17
	CategoryPointObstacle           AircraftCategory = 18
	CategoryClusterObstacle         AircraftCategory = 19
	CategoryLineObstacle            AircraftCategory = 20
)

func (ac AircraftCategory) String() string {
	switch ac {
	case CategoryNoInformation:
		return "No Info"
	case CategoryNoADSB:
		return "No ADSB"
	case CategoryLight:
		return "Light"
	case CategorySmall:
		return "Small"
	case CategoryLarge:
		return "Large"
	case CategoryHighVortexLarge:
		return "High Vortex Large"
	case CategoryHeavy:
		return "Heavy"
	case CategoryHighPerformance:
		return "High Performance"
	case CategoryRotorcraft:
		return "Rotorcraft"
	case CategoryGlider:
		return "Glider"
	case CategoryLighterThanAir:
		return "Lighter than Air"
	case CategoryParachutist:
		return "Parachutist"
	case CategoryUltraLight:
		return "Ultra Light"
	case CategoryReserved:
		return "Reserved"
	case CategoryUnmannedAerialVehicle:
		return "UAV"
	case CategorySpaceVehicle:
		return "Space Vehicle"
	case CategorySurfaceVehicleEmergency:
		return "Surface Vehicle (Emergency)"
	case CategorySurfaceVehicleService:
		return "Surface Vehicle (Service)"
	case CategoryPointObstacle:
		return "Point Obstacle"
	case CategoryClusterObstacle:
		return "Cluster Obstacle"
	case CategoryLineObstacle:
		return "Line Obstacle"
	default:
		return "Unknown"
	}
}

// StateVector represents the state of a single aircraft at a particular time.
type StateVector struct {
	// Icao24: Unique ICAO 24-bit address of the transponder in hex string representation.
	Icao24 string `json:"icao24"`
	// Callsign: Callsign of the vehicle (8 chars). Can be nil if no callsign has been received.
	Callsign *string `json:"callsign,omitempty"`
	// OriginCountry: Country of origin for the given transponder.
	OriginCountry string `json:"origin_country"`
	// TimePosition: Unix timestamp (seconds) for the last position update. Can be nil.
	TimePosition *int64 `json:"time_position,omitempty"`
	// LastContact: Unix timestamp (seconds) for the last update in general.
	LastContact int64 `json:"last_contact"`
	// Longitude: WGS-84 longitude in decimal degrees. Can be nil.
	Longitude *float64 `json:"longitude,omitempty"`
	// Latitude: WGS-84 latitude in decimal degrees. Can be nil.
	Latitude *float64 `json:"latitude,omitempty"`
	// BaroAltitude: Barometric altitude in meters. Can be nil.
	BaroAltitude *float64 `json:"baro_altitude,omitempty"`
	// OnGround: Boolean value which indicates if the position was retrieved from a surface position report.
	OnGround bool `json:"on_ground"`
	// Velocity: Velocity over ground in m/s. Can be nil.
	Velocity *float64 `json:"velocity,omitempty"`
	// TrueTrack: True track in decimal degrees (north is 0°). Can be nil.
	TrueTrack *float64 `json:"true_track,omitempty"`
	// VerticalRate: Vertical rate in m/s. Positive is climbing, negative is descending. Can be nil.
	VerticalRate *float64 `json:"vertical_rate,omitempty"`
	// Sensors: IDs of the receivers which contributed to this state vector. Can be nil.
	Sensors []int `json:"sensors,omitempty"`
	// GeoAltitude: Geometric altitude in meters. Can be nil.
	GeoAltitude *float64 `json:"geo_altitude,omitempty"`
	// Squawk: The transponder code aka Squawk. Can be nil.
	Squawk *string `json:"squawk,omitempty"`
	// SPI: Whether flight status indicates special purpose indicator.
	SPI bool `json:"spi"`
	// PositionSource: Origin of this state's position.
	PositionSource PositionSource `json:"position_source"`
	// Category: The aircraft category.
	Category AircraftCategory `json:"category"`
}

// UnmarshalJSON implements the json.Unmarshaler interface for StateVector.
// The OpenSky API returns state vectors as a JSON array of mixed types.
// This function manually parses the array, which can have 17 or 18 elements.
func (sv *StateVector) UnmarshalJSON(data []byte) error {
	var raw []any
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("failed to unmarshal state vector raw array: %w", err)
	}

	// The API documentation specifies 17 required fields and 1 optional field.
	if len(raw) < 17 {
		return fmt.Errorf("invalid state vector length: got %d, want at least 17", len(raw))
	}

	// Helper function to safely extract and type-assert values from the raw slice.
	// JSON numbers are unmarshaled as float64 into an `any` slice.
	getStr := func(i int) string {
		if val, ok := raw[i].(string); ok {
			return val
		}
		return ""
	}
	getStrPtr := func(i int) *string {
		if val, ok := raw[i].(string); ok {
			return &val
		}
		return nil
	}
	getFloat64Ptr := func(i int) *float64 {
		if val, ok := raw[i].(float64); ok {
			return &val
		}
		return nil
	}
	getInt64 := func(i int) int64 {
		if val, ok := raw[i].(float64); ok {
			return int64(val)
		}
		return 0
	}
	getInt64Ptr := func(i int) *int64 {
		if val, ok := raw[i].(float64); ok {
			ts := int64(val)
			return &ts
		}
		return nil
	}
	getBool := func(i int) bool {
		if val, ok := raw[i].(bool); ok {
			return val
		}
		return false
	}
	getIntSlice := func(i int) []int {
		if raw[i] == nil {
			return nil
		}
		if slice, ok := raw[i].([]any); ok {
			result := make([]int, len(slice))
			for j, item := range slice {
				if val, ok := item.(float64); ok { // JSON numbers are float64
					result[j] = int(val)
				}
			}
			return result
		}
		return nil
	}

	sv.Icao24 = getStr(0)
	sv.Callsign = getStrPtr(1)
	sv.OriginCountry = getStr(2)
	sv.TimePosition = getInt64Ptr(3)
	sv.LastContact = getInt64(4)
	sv.Longitude = getFloat64Ptr(5)
	sv.Latitude = getFloat64Ptr(6)
	sv.BaroAltitude = getFloat64Ptr(7)
	sv.OnGround = getBool(8)
	sv.Velocity = getFloat64Ptr(9)
	sv.TrueTrack = getFloat64Ptr(10)
	sv.VerticalRate = getFloat64Ptr(11)
	sv.Sensors = getIntSlice(12)
	sv.GeoAltitude = getFloat64Ptr(13)
	sv.Squawk = getStrPtr(14)
	sv.SPI = getBool(15)
	sv.PositionSource = PositionSource(getInt64(16))

	// The 'category' field is optional (18th element, index 17).
	// If it's not present, the zero value (CategoryNoInformation) will be used.
	if len(raw) > 17 {
		sv.Category = AircraftCategory(getInt64(17))
	}

	return nil
}
