package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"text/template"
	"time"
)

type StatusData struct {
	WaterStatus struct {
		Water       int    `json:"water"`
		StatusWater string `json:"statusWater"`
	}
	WindStatus struct {
		Wind       int    `json:"wind"`
		StatusWind string `json:"statusWind"`
	}
}

func main() {
	go AutoReloadJSON()
	http.HandleFunc("/", AutoReloadWeb)
	fmt.Println("listening on PORT: ", "8080")
	http.ListenAndServe(":8080", nil)
}

func AutoReloadJSON() {
	for {
		min := 1
		max := 100

		water := rand.Intn(max-min) + 1
		wind := rand.Intn(max-min) + 1

		data := StatusData{}
		data.WaterStatus.Water = water
		switch {
		case water < 5:
			data.WaterStatus.StatusWater = "Aman"
		case water >= 5 && water <= 8:
			data.WaterStatus.StatusWater = "Siaga"
		case water > 8:
			data.WaterStatus.StatusWater = "Bahaya"
		default:
			data.WaterStatus.StatusWater = "Status Tidak Diketahui"
		}

		data.WindStatus.Wind = wind
		switch {
		case wind < 6:
			data.WindStatus.StatusWind = "Aman"
		case wind >= 6 && wind <= 15:
			data.WindStatus.StatusWind = "Siaga"
		case wind > 15:
			data.WindStatus.StatusWind = "Bahaya"
		default:
			data.WindStatus.StatusWind = "Status Tidak Diketahui"
		}

		jsonData, err := json.Marshal(data)

		if err != nil {
			log.Fatal("[error] occured while marshaling status data: ", err.Error())
		}

		if err = ioutil.WriteFile("data/data.json", jsonData, 0644); err != nil {
			log.Fatal("[error] occured while writing json data: ", err.Error())
		}

		time.Sleep(time.Second * 15)
	}
}

func AutoReloadWeb(w http.ResponseWriter, r *http.Request) {
	fileData, err := ioutil.ReadFile("data/data.json")

	if err != nil {
		log.Fatal("[error] error occured while reading data.json file: ", err.Error())
	}

	var statusData StatusData

	err = json.Unmarshal(fileData, &statusData)

	if err != nil {
		log.Fatal("[error] error occured while unmarshaling data.json file: ", err.Error())
	}

	waterValue := statusData.WaterStatus.Water
	waterStatus := statusData.WaterStatus.StatusWater

	windValue := statusData.WindStatus.Wind
	windStatus := statusData.WindStatus.StatusWind

	data := map[string]interface{}{
		"waterValue":  waterValue,
		"waterStatus": waterStatus,
		"windValue":   windValue,
		"windStatus":  windStatus,
	}

	tpl, err := template.ParseFiles("view/index.html")

	if err != nil {
		log.Fatal("[error] error occured while parsing html: ", err.Error())
	}

	tpl.Execute(w, data)
}
