package main

import (
	clients "actions-service/internal/clients"
	"actions-service/internal/config"
	"context"
	"encoding/json"
	"log"
)



func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Couldn't load environment %v", err)
	}
	client, err := clients.NewClientWithResponses(cfg.BackendUrl)
	if err != nil {
		log.Fatalf("Couldn't create client %v", err)
	}
	ctx := context.Background()
	response, err := client.GetApiWorkcenterWithResponse(ctx)
	if err != nil {
		log.Fatalf("Something went wrong calling the backend %v", err)
	}
	if response.HTTPResponse.StatusCode == 200 {
		rawBody := response.Body
		var jsonData []map[string]interface{}
	err := json.Unmarshal(rawBody, &jsonData)
	if err != nil {		
		log.Fatalf("Error deserializing the JSON response: %v", err)
	}

	// Recórrer i modificar camps problemàtics
	for _, workcenter := range jsonData {
		if _, ok := workcenter["createdOn"]; ok {
			workcenter["createdOn"] = nil // Estableix el camp a nil
		}
		if _, ok := workcenter["updatedOn"]; ok {
			workcenter["updatedOn"] = nil // Estableix el camp a nil
		}
	}

	// Torna a serialitzar el JSON modificat
	modifiedBody, err := json.Marshal(jsonData)
	if err != nil {
		log.Fatalf("Error re-serializing the JSON: %v", err)
	}

	// Deserialitza al tipus final
	var workcenters []clients.Workcenter
	err = json.Unmarshal(modifiedBody, &workcenters)
	if err != nil {
		log.Fatalf("Error deserializing into Workcenter: %v", err)
	}

	// Ara pots treballar amb workcenters
	for _, workcenter := range workcenters {
		log.Printf("Workcenter: %+v\n", *workcenter.Description)
		log.Printf("Workcenter: %+v\n", *workcenter.Disabled)
	}
	}
}