package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

//Goal
//----
//You will be entering an HTTP maze. Your goal is to get out of it!
//
//API
//---
//ENDPOINT = "https://territory-pizza-liz-apart.trycloudflare.com"
//The maze will be represented as an API with one GET endpoint.
//To enter the maze, head over to https://territory-pizza-liz-apart.trycloudflare.com
//To go to a specific step in the maze, GET /<STEP_ID>
//The final step of the maze will return a "CONGRATS" message.
//
//Instructions
//------------
//Print the STEP_ID of the final step

var baseURL = "https://territory-pizza-liz-apart.trycloudflare.com"
var stepMap = map[string]bool{}
var nextStepID = ""
var finalStepID = ""
var complete = false

func main() {
	_, err := getNextSteps(nextStepID)
	if err != nil {
		fmt.Println("failed to get next steps, oh no!:", err)
		return
	}
	for !complete {
		for stepID, _ := range stepMap {
			nextStepID = stepID
			nextSteps, err := getNextSteps(nextStepID)
			if err != nil {
				fmt.Println("failed to get next steps, oh no!:", err)
				return
			}
			if len(nextSteps) == 0 {
				fmt.Println("received an empty response for next steps and maze is still incomplete....")
			}
		}
	}

	fmt.Println("FINISHED THE MAZE! FINAL STEP ID:", finalStepID)
}

func insertNextStepsRespToMap(nextSteps []string) {
	for _, stepID := range nextSteps {
		if _, ok := stepMap[stepID]; !ok {
			stepMap[stepID] = true
		}
	}
}

type NextStepsResp struct {
	NextSteps []string `json:"next_steps"`
}

func getNextSteps(stepID string) ([]string, error) {
	fmt.Println("fetching next steps for stepID:", stepID)
	nextSteps := NextStepsResp{}
	client := &http.Client{}

	url := fmt.Sprintf("%s/%s", baseURL, stepID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("failed to build request:", err)
		return nextSteps.NextSteps, err
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println("failed to make request:", err)
		return nextSteps.NextSteps, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusServiceUnavailable {
		fmt.Println("service unavailable, retrying...")
		return getNextSteps(stepID)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("failed to read response body:", err)
		return nextSteps.NextSteps, err
	}

	if strings.Contains(string(body), "CONGRATS") {
		complete = true
		finalStepID = stepID
		return nextSteps.NextSteps, nil
	}

	err = json.Unmarshal(body, &nextSteps)
	if err != nil {
		fmt.Println("failed to unmarshal response body:", err)
		return nextSteps.NextSteps, err
	}

	// add our next steps to the map
	insertNextStepsRespToMap(nextSteps.NextSteps)

	// remove the current stepID from the map since we are done with it
	delete(stepMap, stepID)
	return nextSteps.NextSteps, nil
}
