package register

import (
	"path/filepath"
	"runtime"
	"sort"
)

type registeredDay struct {
	callback  func(string) // Callback function, as defined by the registering function
	inputFile string       // location of the input-file, should be named 'input.txt' and live in the package-directory
}

// All days, registered by the packages
var registeredDays = make(map[string]registeredDay)

func Day(nrDay string, exec func(string)) {
	_, b, _, _ := runtime.Caller(1)
	packageDir := filepath.Dir(b)

	registeredDays[nrDay] = registeredDay{exec, packageDir + "/input.txt"}
}

// Execute the selected puzzle,
// by executing the callback-function with the location of the inputfile as parameter
func ExecDay(nrDay string) {
	var day registeredDay = registeredDays[nrDay]
	day.callback(day.inputFile)
}

func GetAllDays() (nrDays []string) {
	for nrDay := range registeredDays {
		nrDays = append(nrDays, nrDay)
	}

	sort.Strings(nrDays)

	return
}
