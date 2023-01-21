package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
)

func LoadEnvVariables() {

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func apiWarmup() {
	type Input struct {
		Prompt                 string  `json:"prompt"`
		Width                  int     `json:"width"`
		Height                 int     `json:"height"`
		Prompt_Strength        float64 `json:"prompt_strength"`
		Num_Outputs            int     `json:"num_outputs"`
		Num_Interference_Steps int     `json:"num_interference_steps"`
		Guidance_Scale         float64 `json:"guidance_scale"`
	}
	type Post struct {
		Version string `json:"version"`
		Input   Input  `json:"input"`
	}

	posturl := "https://api.replicate.com/v1/predictions"
	values := Post{
		Version: "c24bbf13332c755f9e1c8b3f10c7f438889145def57d554a74ea751dc5e3b509",
		Input: Input{
			Prompt:                 "I like turtles",
			Width:                  128,
			Height:                 128,
			Prompt_Strength:        .5,
			Num_Outputs:            1,
			Num_Interference_Steps: 1,
			Guidance_Scale:         1,
		},
	}
	json_data, err := json.Marshal(values)

	if err != nil {
		log.Fatal(err)
	}

	r, err := http.NewRequest("POST", posturl, bytes.NewBuffer(json_data))
	if err != nil {
		panic(err)
	}
	r.Header.Add("Authorization", "Token "+os.Getenv("NUXT_REP_API_KEY"))
	r.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var res map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&res)

	fmt.Println(res)
}

func runCronJobs() {
	s := gocron.NewScheduler(time.UTC)

	s.Every(10).Minutes().Do(func() {
		apiWarmup()
	})

	s.StartBlocking()
}

func init() {
	LoadEnvVariables()
}

func main() {
	runCronJobs()
}
