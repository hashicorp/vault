package frequency

import (
	"encoding/json"
	"github.com/nbutton23/zxcvbn-go/data"
	"log"
)

type FrequencyList struct {
	Name string
	List []string
}

var FrequencyLists = make(map[string]FrequencyList)

func init() {
	maleFilePath := getAsset("data/MaleNames.json")
	femaleFilePath := getAsset("data/FemaleNames.json")
	surnameFilePath := getAsset("data/Surnames.json")
	englishFilePath := getAsset("data/English.json")
	passwordsFilePath := getAsset("data/Passwords.json")

	FrequencyLists["MaleNames"] = GetStringListFromAsset(maleFilePath, "MaleNames")
	FrequencyLists["FemaleNames"] = GetStringListFromAsset(femaleFilePath, "FemaleNames")
	FrequencyLists["Surname"] = GetStringListFromAsset(surnameFilePath, "Surname")
	FrequencyLists["English"] = GetStringListFromAsset(englishFilePath, "English")
	FrequencyLists["Passwords"] = GetStringListFromAsset(passwordsFilePath, "Passwords")

}
func getAsset(name string) []byte {
	data, err := zxcvbn_data.Asset(name)
	if err != nil {
		panic("Error getting asset " + name)
	}

	return data
}
func GetStringListFromAsset(data []byte, name string) FrequencyList {

	var tempList FrequencyList
	err := json.Unmarshal(data, &tempList)
	if err != nil {
		log.Fatal(err)
	}
	tempList.Name = name
	return tempList
}
