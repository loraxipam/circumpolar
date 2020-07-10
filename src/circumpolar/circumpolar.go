package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

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
func printPairs(pairs []havers2.Coord, r float64, u string, outputJSON bool) {
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
		fmt.Printf("Distances from %-.3f, %-.3f [using a %.1f %s radius]\n", pairs[0].Lat, pairs[0].Lon, r, u)
		for i := 1; i < len(pairs); i++ {
			fmt.Printf(" %-8.3f %-8.3f    %.f %s\t%.f°\n", pairs[i].Lat, pairs[i].Lon, pairs[0].S2Point.Distance(pairs[i].S2Point).Radians()*r, u, -(s2.TurnAngle(pole.S2Point, pairs[0].S2Point, pairs[i].S2Point).Degrees() - 180.0))
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

func main() {
	// The variables used internally
	var (
		outputJSON, kilo, mile bool
		radius                 float64
		unit                   string = "NM"
		locations              []havers2.Coord
		err                    error
	)

	// Get command line flags
	flag.BoolVar(&outputJSON, "json", false, "Output results as JSON")
	flag.BoolVar(&kilo, "kilo", false, "Output station distances in kilometers")
	flag.BoolVar(&mile, "mile", false, "Output station distances in statue miles")
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
		flag.Usage()
		os.Exit(1)
	}

	// Create lat/lon pairs
	// FYI, use flag.Args here instead of os.Args because flag.Args already has the cruft removed
	locations, err = makePairs(flag.Args())
	if err != nil {
		log.Println(err)
	}

	// Printout the results
	printPairs(locations, radius, unit, outputJSON)
}
