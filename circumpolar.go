package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"net/http"
	"crypto/tls"
	"io/ioutil"

	json "github.com/json-iterator/go"

	"github.com/golang/geo/s2"
	"github.com/loraxipam/havers2"
)

// jpair is an overloaded havers2.Coord object used for displaying JSON
type jpair struct {
	Coord havers2.Coord `json:"coord,omitempty"`
	Dist  float64       `json:"distance"`
	Head  float64       `json:"heading,omitempty"`
	Index int           `json:"index"`
}

// {
//   "result": [
//     {
//       "date": 2020.7897,
//       "elevation": 0,
//       "declination": -6.88502,
//       "latitude": 29.13,
//       "declnation_sv": -0.07518,
//       "declination_uncertainty": 0.34714,
//       "longitude": -80.96
//     }
//   ],
//   "model": "WMM-2020",
//   "units": {
//     "elevation": "km",
//     "declination": "degrees",
//     "declination_sv": "degrees",
//     "latitude": "degrees",
//     "declination_uncertainty": "degrees",
//     "longitude": "degrees"
//   },
//   "version": "0.5.1.11"
// }

type MagDec struct {
	Date float64 `json:"date"`
	Elevation float64 `json:"elevation"`
	Declination float64 `json:"declination"`
	Latitude float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	DeclSv float64 `json:"declnation_sv"`
	DeclUn float64 `json:"declination_uncertainty"`
}

// makePairs turns the command line arguments into an array of coordinates
func makePairs(args []string) (spread []havers2.Coord, err error) {
	var tmpPnts = make([][2]float64, len(args)/2)

	// Process args
	for plen := range tmpPnts {
		var tmpPnt [2]float64
		tmpPnt[0], err = strconv.ParseFloat(args[plen*2], 32)
		if err != nil {
			return nil, err
		}
		tmpPnt[1], err = strconv.ParseFloat(args[plen*2+1], 32)
		if err != nil {
			return nil, err
		}
		tmpPnts[plen] = tmpPnt
	}

	for k := 0; k < len(tmpPnts); k++ {
		tmp := havers2.Coord{Lat: tmpPnts[k][0], Lon: tmpPnts[k][1]}
		tmp.Calc()
		spread = append(spread, tmp)
	}

	return spread, err
}

// printPairs sends results to stdout in either text columns or JSON
func printPairs(pairs []havers2.Coord, g, r float64, u string, outputJSON bool) {
	pole := havers2.Coord{Lat: 90.0, Lon: 0.0}
	pole.Calc()

	if outputJSON {
		// For JSON, we'll make an array with some extras thrown in.
		jpairs := make([]jpair, len(pairs))

		// Populate the JSON-producer array
		for key, val := range pairs {
			jpairs[key].Index = key
			jpairs[key].Coord = val
			jpairs[key].Dist = pairs[0].S2Point.Distance(val.S2Point).Radians() * r
			jpairs[key].Head = -(s2.TurnAngle(pole.S2Point, pairs[0].S2Point, val.S2Point).Degrees() - 180.0)
		}

		j, err := json.Marshal(jpairs)
		if err != nil {
			log.Println("Cannot marshal the JSON:", err)
		} else {
			fmt.Printf(string(j))
		}
	} else {
		fmt.Printf("Distances from %-.3f, %-.3f [using a %.1f %s radius. Magnetic declination there is %-.2f]\n", pairs[0].Lat, pairs[0].Lon, r, u, g)
		for i := 1; i < len(pairs); i++ {
			angle := -(s2.TurnAngle(pole.S2Point, pairs[0].S2Point, pairs[i].S2Point).Degrees() - 180.0)
			fmt.Printf(" %-8.3f %-8.3f    %.f %s\t%.f°\t[%.f°]\n", pairs[i].Lat, pairs[i].Lon, pairs[0].S2Point.Distance(pairs[i].S2Point).Radians()*r, u, angle, angle + g)
		}
	}
}

// contains tells whether "a" contains "x".
func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

// get compass declination from the web
func getDeclinationInfoFromNOAA(c havers2.Coord) (decl float64, err error) {

	// We need a TLS session
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// We need a client for the TLS session
	client := &http.Client{Transport: transport}

	// We need a declination service URL for NOAA's API
	urlTemplate := "https://www.ngdc.noaa.gov/geomag-web/calculators/calculateDeclination?lat1=%f&lon1=%f&resultFormat=json&model=WMM&magneticComponent=d"

	apiURL := fmt.Sprintf(urlTemplate,c.Lat,c.Lon)

	// Make the call
	responseBody, err := client.Get(apiURL)
	if err != nil {
		return -999.99, err
	}

	// Close the session once we're done
	defer responseBody.Body.Close()

	// Now parse the result
	apiResponse, err := ioutil.ReadAll(responseBody.Body)
	if err != nil {
		return -999.99, err
	}

	var returned map[string][]MagDec

	// Unmarshal the JSON
	err = json.Unmarshal(apiResponse, &returned)

	decl = returned["result"][0].Declination

	return decl, nil
}

func main() {
	// The variables used internally
	var (
		outputJSON, kilo, mile, home bool
		radius, geoDecl              float64
		unit                         string = "NM"
		locations                    []havers2.Coord
		err                          error
	)

	// Get command line flags
	flag.BoolVar(&outputJSON, "json", false, "Output results as JSON")
	flag.BoolVar(&kilo, "kilo", false, "Output station distances in kilometers")
	flag.BoolVar(&mile, "mile", false, "Output station distances in statue miles")
	flag.BoolVar(&home, "home", false, "Stay home. Don't query NOAA for declination")
	flag.Float64Var(&radius, "radius", havers2.EarthRadiusNM, "Assign the sphere's radius to this value instead of Earth's nautical miles")
	flag.Parse()

	// Set the radius and units based on the flags
	switch {
	case kilo:
		unit = "km"
		if !contains(os.Args, "-radius") {
			radius = havers2.EarthRadiusKm
		}
	case mile:
		unit = "mi"
		if !contains(os.Args, "-radius") {
			radius = havers2.EarthRadiusMi
		}
	}

	// Did they pass ANYTHING?
	if len(flag.Args()) < 2 {
		fmt.Println("\ncircumpolar latA lonA latX lonX [latY lonY latZ lonZ...]")
		fmt.Println("    where lat/lon values are decimal with negative S and W values")
		flag.Usage()
		os.Exit(1)
	}

	// Create lat/lon pairs
	// FYI, use flag.Args here instead of os.Args because flag.Args already has the cruft removed
	locations, err = makePairs(flag.Args())
	if err != nil {
		log.Println(err)
	}

	if home {
		geoDecl = 0
	} else {
		geoDecl, err = getDeclinationInfoFromNOAA(locations[0])
	}

	// Printout the results
	printPairs(locations, geoDecl, radius, unit, outputJSON)

}
