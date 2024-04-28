package reader

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

func OpenFile(inputname string) (io.Reader, error) {
	if strings.HasSuffix(inputname, ".zip") {
		zfc, err := zip.OpenReader(inputname)
		if err != nil {
			return nil, err
		}
		// TODO: Do something better with closer.
		// defer zfc.Close()

		for _, f := range zfc.Reader.File {
			if strings.HasSuffix(f.Name, "/Records.json") {
				return f.Open()
			}
		}
		return nil, fmt.Errorf("could not find Records.json inside zip file")
	}
	if strings.HasSuffix(inputname, ".json") {
		return os.Open(inputname)
	}
	return nil, fmt.Errorf("only .zip and .json are supported")
}

// DecodeJson attempts to read takeout-compatible JSON from the given reader.
func DecodeJson(reader io.Reader) ([]Location, error) {
	decoder := json.NewDecoder(reader)

	// Read the following opening tokens:
	// 1. open brace '{'
	// 2. "locations" field name,
	// 3. the array value's opening bracket '['
	for i := 0; i < 3; i++ {
		_, err := decoder.Token()
		if err != nil {
			return nil, fmt.Errorf("decoding opening token: %v", err)
		}
	}

	var locs []Location
	for decoder.More() {
		loc := Location{}
		err := decoder.Decode(&loc)
		if err != nil {
			return nil, err
		}
		locs = append(locs, loc)
	}
	return locs, nil
}

type Location struct {
	Timestamp        string       `json:"timestamp"`
	LatitudeE7       int          `json:"latitudeE7"`
	LongitudeE7      int          `json:"longitudeE7"`
	Accuracy         int          `json:"accuracy"`
	Altitude         int          `json:"altitude,omitempty"`
	VerticalAccuracy int          `json:"verticalAccuracy,omitempty"`
	Activity         []activities `json:"activity,omitempty"`
	Velocity         int          `json:"velocity,omitempty"`
	Heading          int          `json:"heading,omitempty"`

	// Maybe useful?
	FormFactor      string `json:"formFactor"` // PHONE
	BatteryCharging bool   `json:"batteryCharging"`
	Source          string `json:"source"`       // WIFI, GPS
	PlatformType    string `json:"platformType"` // ANDROID
	// locationMetadata has an array of wifi scans.
}

func (l Location) ParsedTimestamp() (time.Time, error) {
	return time.Parse(time.RFC3339, l.Timestamp)
}

type activities struct {
	Timestamp string     `json:"timestamp"`
	Activity  []activity `json:"activity"`
}

type activity struct {
	Type       string `json:"type"`
	Confidence int    `json:"confidence"`
}
