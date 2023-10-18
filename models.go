package opensky

type State struct {
	ICAO24             string  `json:"icao24"`
	Callsign           string  `json:"callsign"`
	OriginCountry      string  `json:"origin_country"`
	TimePosition       int     `json:"time_position"`
	LastContact        float64 `json:"last_contact"`
	Longitude          float32 `json:"longitude"`
	Latitude           float32 `json:"latitude"`
	BarometricAltitude float32 `json:"baro_altitude"`
	OnGround           bool    `json:"on_ground"`
	Velocity           float32 `json:"velocity"`
	TrueTrack          float32 `json:"true_track"`
	VerticalRate       float32 `json:"vertical_rate"`
	Sensors            []int   `json:"sensors"`
	GeometricAltitude  float32 `json:"geo_altitude"`
	Squawk             string  `json:"squawk"`
	SPI                bool    `json:"spi"`
	PositionSource     float64 `json:"position_source"`
	Category           float64 `json:"category"`
}

type StateResponse struct {
	Time   int     `json:"time"`
	States []State `json:"states"`
}

type rawStateResponse struct {
	Time   int             `json:"time"`
	States [][]interface{} `json:"states"`
}

func (raw *rawStateResponse) parse() StateResponse {
	resp := StateResponse{}
	resp.Time = raw.Time

	for _, s := range raw.States {
		parsed, err := parseStateRow(s)
		if err != nil {
			continue
		}
		resp.States = append(resp.States, parsed)
	}

	return resp
}

func parseStateRow(rowdata []interface{}) (State, error) {
	state := State{
		ICAO24:         rowdata[0].(string),
		OriginCountry:  rowdata[2].(string),
		LastContact:    rowdata[4].(float64),
		OnGround:       rowdata[8].(bool),
		SPI:            rowdata[15].(bool),
		PositionSource: rowdata[16].(float64),
		//Category:       rowdata[17].(float64),
	}

	if len(rowdata) == 18 {
		state.Category = rowdata[17].(float64)
	}

	return state, nil
}
