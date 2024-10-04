package  main 

import (
	"fmt"
	"net/http"
	"io"
	"encoding/json"
	"github.com/fatih/color"
	"time"
	"os"
)

type Weather struct {
	Location struct {
		Name           string 	`json:"name"`
		Region         string 	`json:"region"`
		Country        string	`json:"country"`
	} `json:"location"`
	Current struct {
		TempC            float64 `json:"temp_c"`
		Condition        struct {
			Text string `json:"text`
		} `json:"condition"`
		Uv         float64 `json:"uv"`
	} `json:"current"`
	Forecast struct {
		Forecastday []struct {
			Date      string `json:"date"`
			Day       struct {
				MaxtempC          float64 `json:"maxtemp_c"`
				MintempC          float64 `json:"mintemp_c"`
				AvgtempC          float64 `json:"avgtemp_c"`
				DailyChanceOfRain int     `json:"daily_chance_of_rain"`
				DailyChanceOfSnow int     `json:"daily_chance_of_snow"`
				Condition         struct {
					Text string `json:"text"`
				} `json:"condition"`
				Uv float64 `json:"uv"`
			} `json:"day"`
			Hour []struct {
				TimeEpoch int64     `json:"time_epoch"`
				TempC     float64 `json:"temp_c"`
				Condition struct {
					Text string `json:"text"`
				} `json:"condition"`

				FeelslikeC   float64 `json:"feelslike_c"`
				ChanceOfRain int     `json:"chance_of_rain"`
				ChanceOfSnow int     `json:"chance_of_snow"`
				Uv           float64     `json:"uv"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func main(){
	q := "Bangalore"
	if len(os.Args) >= 2{
		q = os.Args[1]
	}
	apiKey := "replace with api key"
	req :=  fmt.Sprintf("https://api.weatherapi.com/v1/forecast.json?key=%s&q=%s&aqi=no&days=1&tp=1440", apiKey, q)
	res, err := http.Get(req) 
	if err != nil{
		fmt.Println("something went wrong!")
		panic(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		fmt.Println("theres somethign wrong with this api!")
		panic(err)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("error reading content")
		panic(err)
	}
	var weather Weather 
	err = json.Unmarshal(body, &weather)
	if err != nil{
		fmt.Println("error decoding json")
		panic(err)
	}
	fmt.Println(weather.Location.Name, weather.Location.Region, weather.Location.Country) 
	fmt.Println("Date ",weather.Forecast.Forecastday[0].Date)
	fmt.Println("Temp ",weather.Current.TempC)
	for _, i := range weather.Forecast.Forecastday[0].Hour{
		ti := time.Unix(i.TimeEpoch, 0)
		if ti.Before(time.Now()){
			continue
		}
		formatted := ti.Format("15:03")
		if i.ChanceOfRain > 50 {
			color.Cyan("%s %s %f",formatted,i.Condition.Text, i.TempC )
		}else if i.ChanceOfSnow > 50 {
			color.White("%s %s %f",formatted,i.Condition.Text, i.TempC )
		}else if i.TempC > 40 {
			color.Red("%s %s %f",formatted,i.Condition.Text, i.TempC )
            	}else if i.TempC > 20 {
			color.Yellow("%s %s %f",formatted,i.Condition.Text, i.TempC )
		}else{
			color.Green("%s %s %f",formatted,i.Condition.Text, i.TempC )
		}
	}

}
