package main

import (
	"fmt"
	"log"
	"os"

	growthbook "github.com/growthbook/growthbook-golang"
)

func main() {
	// Set up development logger: logs all messages from GrowthBook SDK
	// and exits on errors.
	// growthbook.SetLogger(&growthbook.DevLogger{})

	// Read JSON feature file.
	featureJSON, err := os.ReadFile("features.json")
	if err != nil {
		log.Fatal(err)
	}

	// Parse feature map from JSON.
	features := growthbook.ParseFeatureMap(featureJSON)

	// Create context and main GrowthBook object.
	context := growthbook.NewContext().
		WithFeatures(features).
		WithAttributes(growthbook.Attributes{
			"country": "US",
			"browser": "firefox",
		})
	gb := growthbook.New(context)

	// Perform feature test.
	fmt.Print("test-feature (US, firefox): ")
	if gb.Feature("test-feature").On {
		fmt.Print("ON")
	} else {
		fmt.Println("OFF")
	}
	fmt.Println("  value =", gb.Feature("test-feature").Value)

	// Perform feature test with different user attributes.
	gb.WithAttributes(growthbook.Attributes{
		"country": "AT",
		"browser": "firefox",
	})
	fmt.Print("test-feature (AT, firefox): ")
	if gb.Feature("test-feature").On {
		fmt.Print("ON")
	} else {
		fmt.Println("OFF")
	}
	fmt.Println("  value =", gb.Feature("test-feature").Value)

	// Feature value lookup with default.
	color := gb.Feature("signup-button-color").GetValueWithDefault("blue")
	fmt.Printf("\nsignup-button-color: %s\n", color)

	// Run simple experiment.
	experiment :=
		growthbook.NewExperiment("my-experiment").
			WithVariations("A", "B")

	result := gb.Run(experiment)

	fmt.Printf("\nmy-experiment value: %s\n\n", result.Value)

	// Run more complex experiment.
	experiment2 :=
		growthbook.NewExperiment("complex-experiment").
			WithVariations(
				map[string]string{"color": "blue", "size": "small"},
				map[string]string{"color": "green", "size": "large"},
			).
			WithWeights(0.8, 0.2).
			WithCoverage(0.5)

	result2 := gb.Run(experiment2)
	fmt.Println("complex-experiment values:",
		result2.Value.(map[string]string)["color"],
		result2.Value.(map[string]string)["size"])
}
