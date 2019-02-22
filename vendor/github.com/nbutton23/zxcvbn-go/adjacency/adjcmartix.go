package adjacency

import (
	"encoding/json"
	"log"
	//	"fmt"
	"github.com/nbutton23/zxcvbn-go/data"
)

type AdjacencyGraph struct {
	Graph         map[string][]string
	averageDegree float64
	Name          string
}

var AdjacencyGph = make(map[string]AdjacencyGraph)

func init() {
	AdjacencyGph["qwerty"] = BuildQwerty()
	AdjacencyGph["dvorak"] = BuildDvorak()
	AdjacencyGph["keypad"] = BuildKeypad()
	AdjacencyGph["macKeypad"] = BuildMacKeypad()
	AdjacencyGph["l33t"] = BuildLeet()
}

func BuildQwerty() AdjacencyGraph {
	data, err := zxcvbn_data.Asset("data/Qwerty.json")
	if err != nil {
		panic("Can't find asset")
	}
	return GetAdjancencyGraphFromFile(data, "qwerty")
}
func BuildDvorak() AdjacencyGraph {
	data, err := zxcvbn_data.Asset("data/Dvorak.json")
	if err != nil {
		panic("Can't find asset")
	}
	return GetAdjancencyGraphFromFile(data, "dvorak")
}
func BuildKeypad() AdjacencyGraph {
	data, err := zxcvbn_data.Asset("data/Keypad.json")
	if err != nil {
		panic("Can't find asset")
	}
	return GetAdjancencyGraphFromFile(data, "keypad")
}
func BuildMacKeypad() AdjacencyGraph {
	data, err := zxcvbn_data.Asset("data/MacKeypad.json")
	if err != nil {
		panic("Can't find asset")
	}
	return GetAdjancencyGraphFromFile(data, "mac_keypad")
}
func BuildLeet() AdjacencyGraph {
	data, err := zxcvbn_data.Asset("data/L33t.json")
	if err != nil {
		panic("Can't find asset")
	}
	return GetAdjancencyGraphFromFile(data, "keypad")
}

func GetAdjancencyGraphFromFile(data []byte, name string) AdjacencyGraph {

	var graph AdjacencyGraph
	err := json.Unmarshal(data, &graph)
	if err != nil {
		log.Fatal(err)
	}
	graph.Name = name
	return graph
}

//on qwerty, 'g' has degree 6, being adjacent to 'ftyhbv'. '\' has degree 1.
//this calculates the average over all keys.
//TODO double check that i ported this correctly scoring.coffee ln 5
func (adjGrp AdjacencyGraph) CalculateAvgDegree() float64 {
	if adjGrp.averageDegree != float64(0) {
		return adjGrp.averageDegree
	}
	var avg float64
	var count float64
	for _, value := range adjGrp.Graph {

		for _, char := range value {
			if char != "" || char != " " {
				avg += float64(len(char))
				count++
			}
		}

	}

	adjGrp.averageDegree = avg / count

	return adjGrp.averageDegree
}
