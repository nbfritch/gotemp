package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	invalidHwmonName  = "INVALID"
	invalidHwmonInput = "INVALID"
	hwmonNameGlob     = "/sys/class/hwmon/hwmon*/name"
)

var (
	hwmonName  = flag.String("name", invalidHwmonName, "Name of the hwmon node to measure")
	hwmonInput = flag.String("input", invalidHwmonInput, "Optional: Name of input. ex: temp1_input or blank to average")
	warnTemp   = flag.Uint("warn", 70, "Temp in celsius after which warn class should be set")
	critTemp   = flag.Uint("crit", 90, "Temp in celsius after which crit class should be set")
	verbose    = flag.Bool("verbose", false, "Whether to use verbose logging")
)

type WaybarTemperatureInput struct {
	Text  string `json:"text"`
	Class string `json:"class"`
}

func findFirstHwmonWithName(name string) (string, error) {
	hwmons, err := filepath.Glob(hwmonNameGlob)
	if err != nil {
		return "", err
	}
	for _, f := range hwmons {
		nodeNameBytes, err := os.ReadFile(f)
		if err != nil {
			return "", err
		}

		nameStr := strings.Trim(string(nodeNameBytes), "\n")
		if nameStr == name {
			return strings.ReplaceAll(f, "name", ""), nil
		}
	}

	return "", errors.New("Could not find any matching nodes")
}

func readHwmonInputs(hwmonPath string, input string) (uint64, error) {
	specificInput := "*"
	if len(input) > 0 && input != invalidHwmonInput {
		specificInput = input
	}
	allInputsGlob := string(fmt.Sprintf("%stemp%s_input", hwmonPath, specificInput))
	inputs, err := filepath.Glob(allInputsGlob)
	if err != nil {
		return 0, err
	}

	totalReading := uint64(0)
	totalCount := uint64(0)
	for _, hInput := range inputs {
		reading, err := os.ReadFile(hInput)
		if err != nil {
			return 0, err
		}

		rString := strings.Trim(string(reading), "\n")
		value, err := strconv.ParseUint(rString, 10, 64)
		if err != nil {
			return 0, err
		}

		totalReading = totalReading + value
		totalCount = totalCount + 1
	}

	if totalCount == 0 {
		return 0, errors.New("Error, no inputs found")
	}

	return (totalReading / totalCount), nil
}

func displayTemp(w WaybarTemperatureInput) error {
	o, err := json.Marshal(w)
	if err != nil {
		return err
	}

	fmt.Println(string(o))
	return nil
}

func failError() {
	err := displayTemp(WaybarTemperatureInput{
		Text:  "ERR",
		Class: "critical",
	})

	if err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()

	if *hwmonName == invalidHwmonName {
		if *verbose {

			fmt.Println("Error: -name flag is required")
		}
		failError()
		os.Exit(1)
	}

	hwmonPath, err := findFirstHwmonWithName(*hwmonName)
	if err != nil {
		if *verbose {
			fmt.Printf("Error: hwmon node with name %s not found\n", *hwmonName)
		}
		failError()
		os.Exit(1)
	}

	reading, err := readHwmonInputs(hwmonPath, *hwmonInput)
	if err != nil {
		if *verbose {
			fmt.Printf("%s\n", err)
		}
		os.Exit(1)
	}

	class := "normal"
	actualReading := uint(reading / 1000)
	if actualReading > *critTemp {
		class = "critical"
	} else if actualReading > *warnTemp {
		class = "warning"
	}

	err = displayTemp(WaybarTemperatureInput{
		Text:  strconv.FormatUint(uint64(actualReading), 10),
		Class: class,
	})

	if err != nil {
		if *verbose {
			fmt.Printf("Json serialization error: %s\n", err)
		}
		failError()
		os.Exit(1)
	}
}
