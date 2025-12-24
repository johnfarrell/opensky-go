package opensky

import (
	"encoding/json"
	"reflect"
	"testing"
)

// TestStateVector_UnmarshalJSON uses a table-driven approach to test various
// unmarshaling scenarios for the StateVector struct.
func TestStateVector_UnmarshalJSON(t *testing.T) {
	// Helper functions to create pointers for the 'want' structs
	stringPtr := func(s string) *string { return &s }
	int64Ptr := func(i int64) *int64 { return &i }
	float64Ptr := func(f float64) *float64 { return &f }

	testCases := []struct {
		name    string
		input   string
		want    StateVector
		wantErr bool
	}{
		{
			name:  "Full Vector (18 fields)",
			input: `["a8a776","KLM123  ","Netherlands",1672531199,1672531200,4.764,52.308,11000.45,false,250.7,345.5,-3.25,[1,2,3],11300.22,"1234",false,0,6]`,
			want: StateVector{
				Icao24:         "a8a776",
				Callsign:       stringPtr("KLM123  "),
				OriginCountry:  "Netherlands",
				TimePosition:   int64Ptr(1672531199),
				LastContact:    1672531200,
				Longitude:      float64Ptr(4.764),
				Latitude:       float64Ptr(52.308),
				BaroAltitude:   float64Ptr(11000.45),
				OnGround:       false,
				Velocity:       float64Ptr(250.7),
				TrueTrack:      float64Ptr(345.5),
				VerticalRate:   float64Ptr(-3.25),
				Sensors:        []int{1, 2, 3},
				GeoAltitude:    float64Ptr(11300.22),
				Squawk:         stringPtr("1234"),
				SPI:            false,
				PositionSource: PositionSourceADS_B,
				Category:       CategoryHeavy,
			},
			wantErr: false,
		},
		{
			name:  "Legacy Vector (17 fields, no category)",
			input: `["4b1815","SWR456  ","Switzerland",1672531198,1672531198,8.54,47.45,9500,false,220.1,120.9,0,null,9800,"5678",false,2]`,
			want: StateVector{
				Icao24:         "4b1815",
				Callsign:       stringPtr("SWR456  "),
				OriginCountry:  "Switzerland",
				TimePosition:   int64Ptr(1672531198),
				LastContact:    1672531198,
				Longitude:      float64Ptr(8.54),
				Latitude:       float64Ptr(47.45),
				BaroAltitude:   float64Ptr(9500),
				OnGround:       false,
				Velocity:       float64Ptr(220.1),
				TrueTrack:      float64Ptr(120.9),
				VerticalRate:   float64Ptr(0),
				Sensors:        nil,
				GeoAltitude:    float64Ptr(9800),
				Squawk:         stringPtr("5678"),
				SPI:            false,
				PositionSource: PositionSourceMLAT,
				Category:       CategoryNoInformation, // zero value
			},
			wantErr: false,
		},
		{
			name:  "Vector with all nullable fields as null",
			input: `["c0ffee",null,"United States",null,1672531195,null,null,null,true,null,null,null,null,null,null,false,1,17]`,
			want: StateVector{
				Icao24:         "c0ffee",
				Callsign:       nil,
				OriginCountry:  "United States",
				TimePosition:   nil,
				LastContact:    1672531195,
				Longitude:      nil,
				Latitude:       nil,
				BaroAltitude:   nil,
				OnGround:       true,
				Velocity:       nil,
				TrueTrack:      nil,
				VerticalRate:   nil,
				Sensors:        nil,
				GeoAltitude:    nil,
				Squawk:         nil,
				SPI:            false,
				PositionSource: PositionSourceASTERIX,
				Category:       CategorySurfaceVehicleService,
			},
			wantErr: false,
		},
		{
			name:  "Vector with empty sensor array",
			input: `["c0ffee",null,"United States",null,1672531195,null,null,null,true,null,null,null,[],null,null,false,1,17]`,
			want: StateVector{
				Icao24:         "c0ffee",
				OriginCountry:  "United States",
				LastContact:    1672531195,
				OnGround:       true,
				Sensors:        []int{}, // Note: an empty slice, not nil
				PositionSource: PositionSourceASTERIX,
				Category:       CategorySurfaceVehicleService,
			},
			wantErr: false,
		},
		{
			name:    "Error on vector too short",
			input:   `["a8a776","KLM123"]`,
			want:    StateVector{},
			wantErr: true,
		},
		{
			name:    "Error on non-array JSON",
			input:   `{"icao24": "a8a776"}`,
			want:    StateVector{},
			wantErr: true,
		},
		{
			name: "Handles type mismatch gracefully",
			// The longitude at index 5 is a string, not a number.
			// The custom unmarshaler's getFloat64Ptr should handle this by returning nil.
			input: `["a8a776", "KLM123", "Netherlands", 1672531199, 1672531200, "not-a-number", 52.308, 11000.45, false, 250.7, 345.5, 0, null, 11300.22, "1234", false, 0, 6]`,
			want: StateVector{
				Icao24:         "a8a776",
				Callsign:       stringPtr("KLM123"),
				OriginCountry:  "Netherlands",
				TimePosition:   int64Ptr(1672531199),
				LastContact:    1672531200,
				Longitude:      nil, // This is the key check for this test case
				Latitude:       float64Ptr(52.308),
				BaroAltitude:   float64Ptr(11000.45),
				OnGround:       false,
				Velocity:       float64Ptr(250.7),
				TrueTrack:      float64Ptr(345.5),
				VerticalRate:   float64Ptr(0),
				Sensors:        nil,
				GeoAltitude:    float64Ptr(11300.22),
				Squawk:         stringPtr("1234"),
				SPI:            false,
				PositionSource: PositionSourceADS_B,
				Category:       CategoryHeavy,
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var got StateVector
			err := json.Unmarshal([]byte(tc.input), &got)

			if (err != nil) != tc.wantErr {
				t.Fatalf("Unmarshal() error = %v, wantErr %v", err, tc.wantErr)
			}

			if !tc.wantErr && !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Unmarshal() got = %+v, want %+v", got, tc.want)
			}
		})
	}
}

// TestStateVector_MarshalJSON verifies that a StateVector struct marshals
// into a well-formed JSON object with snake_case and omitempty tags.
func TestStateVector_MarshalJSON(t *testing.T) {
	float64Ptr := func(f float64) *float64 { return &f }

	// A struct with some fields populated and others nil/zero to test omitempty
	sv := StateVector{
		Icao24:         "a8a776",
		Callsign:       nil, // Should be omitted
		OriginCountry:  "Netherlands",
		LastContact:    1672531200,
		Longitude:      float64Ptr(4.764),
		OnGround:       false,
		Velocity:       nil, // Should be omitted
		SPI:            true,
		PositionSource: PositionSourceADS_B,
		Category:       CategoryHeavy,
	}

	// Marshal the struct to JSON bytes
	gotBytes, err := json.Marshal(sv)
	if err != nil {
		t.Fatalf("json.Marshal() failed: %v", err)
	}

	// We can't compare JSON strings directly due to key ordering.
	// The robust way is to unmarshal both into maps and compare the maps.
	var gotMap map[string]any
	if err := json.Unmarshal(gotBytes, &gotMap); err != nil {
		t.Fatalf("Failed to unmarshal result for comparison: %v", err)
	}

	// The 'callsign' and 'velocity' keys should not exist due to omitempty.
	expectedJSON := `{
		"icao24": "a8a776",
		"origin_country": "Netherlands",
		"last_contact": 1672531200,
		"longitude": 4.764,
		"on_ground": false,
		"spi": true,
		"position_source": 0,
		"category": 6
	}`
	var wantMap map[string]any
	if err := json.Unmarshal([]byte(expectedJSON), &wantMap); err != nil {
		t.Fatalf("Failed to unmarshal expected JSON for comparison: %v", err)
	}

	if !reflect.DeepEqual(gotMap, wantMap) {
		t.Errorf("Marshal() result mismatch:\ngot = %s\nwant based on = %s", string(gotBytes), expectedJSON)
	}
}

// TestStateResponse_UnmarshalJSON verifies that the top-level response struct
// correctly unmarshals and uses the custom unmarshaler for its nested slice.
func TestStateResponse_UnmarshalJSON(t *testing.T) {
	input := `{
		"time": 1672531200,
		"states": [
			["a8a776", null, "Netherlands", null, 1672531200, null, null, null, false, null, null, null, null, null, null, false, 0, 6]
		]
	}`

	var response StateResponse
	err := json.Unmarshal([]byte(input), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal StateResponse: %v", err)
	}

	if response.Time != 1672531200 {
		t.Errorf("Incorrect Time field: got %d, want %d", response.Time, 1672531200)
	}

	if len(response.States) != 1 {
		t.Fatalf("Incorrect number of states: got %d, want 1", len(response.States))
	}

	// Verify that the custom unmarshaler was used correctly
	expectedIcao24 := "a8a776"
	if response.States[0].Icao24 != expectedIcao24 {
		t.Errorf("Incorrect Icao24 in nested state: got %s, want %s", response.States[0].Icao24, expectedIcao24)
	}

	expectedCategory := CategoryHeavy
	if response.States[0].Category != expectedCategory {
		t.Errorf("Incorrect Category in nested state: got %v, want %v", response.States[0].Category, expectedCategory)
	}
}
