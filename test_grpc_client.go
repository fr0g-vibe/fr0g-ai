package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip/internal/grpc/pb"
)

func main() {
	// Connect to the gRPC server
	conn, err := grpc.Dial("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := pb.NewPersonaServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Println("=== gRPC Client Test for fr0g-ai-aip ===")
	fmt.Println()

	// Test 1: Create a persona
	fmt.Println("Test 1: Creating a persona...")
	createReq := &pb.CreatePersonaRequest{
		Persona: &pb.Persona{
			Name:   "gRPC Test Persona",
			Topic:  "gRPC Testing",
			Prompt: "You are a test persona created via gRPC client",
			Context: map[string]string{
				"test_type": "grpc_client",
				"timestamp": time.Now().Format(time.RFC3339),
			},
			Rag: []string{"grpc_test_document"},
		},
	}

	createResp, err := client.CreatePersona(ctx, createReq)
	if err != nil {
		log.Printf("Failed to create persona: %v", err)
		return
	}

	fmt.Printf("✓ Created persona with ID: %s\n", createResp.Persona.Id)
	personaID := createResp.Persona.Id

	// Test 2: Get the persona
	fmt.Println("\nTest 2: Retrieving the persona...")
	getReq := &pb.GetPersonaRequest{
		Id: personaID,
	}

	getResp, err := client.GetPersona(ctx, getReq)
	if err != nil {
		log.Printf("Failed to get persona: %v", err)
	} else {
		fmt.Printf("✓ Retrieved persona: %s\n", getResp.Persona.Name)
		fmt.Printf("  Topic: %s\n", getResp.Persona.Topic)
		fmt.Printf("  Context entries: %d\n", len(getResp.Persona.Context))
	}

	// Test 3: Update the persona
	fmt.Println("\nTest 3: Updating the persona...")
	updateReq := &pb.UpdatePersonaRequest{
		Persona: &pb.Persona{
			Id:     personaID,
			Name:   "Updated gRPC Test Persona",
			Topic:  "Updated gRPC Testing",
			Prompt: "You are an updated test persona via gRPC client",
			Context: map[string]string{
				"test_type": "grpc_client",
				"timestamp": time.Now().Format(time.RFC3339),
				"updated":   "true",
			},
			Rag: []string{"updated_grpc_document"},
		},
	}

	updateResp, err := client.UpdatePersona(ctx, updateReq)
	if err != nil {
		log.Printf("Failed to update persona: %v", err)
	} else {
		fmt.Printf("✓ Updated persona: %s\n", updateResp.Persona.Name)
	}

	// Test 4: List personas
	fmt.Println("\nTest 4: Listing all personas...")
	listReq := &pb.ListPersonasRequest{}

	listResp, err := client.ListPersonas(ctx, listReq)
	if err != nil {
		log.Printf("Failed to list personas: %v", err)
	} else {
		fmt.Printf("✓ Found %d personas in total\n", len(listResp.Personas))
		for i, persona := range listResp.Personas {
			if i < 5 { // Show first 5
				fmt.Printf("  - %s (ID: %s)\n", persona.Name, persona.Id)
			}
		}
		if len(listResp.Personas) > 5 {
			fmt.Printf("  ... and %d more\n", len(listResp.Personas)-5)
		}
	}

	// Test 5: Delete the test persona
	fmt.Println("\nTest 5: Deleting the test persona...")
	deleteReq := &pb.DeletePersonaRequest{
		Id: personaID,
	}

	_, err = client.DeletePersona(ctx, deleteReq)
	if err != nil {
		log.Printf("Failed to delete persona: %v", err)
	} else {
		fmt.Printf("✓ Deleted persona with ID: %s\n", personaID)
	}

	// Test 6: Verify deletion
	fmt.Println("\nTest 6: Verifying deletion...")
	_, err = client.GetPersona(ctx, &pb.GetPersonaRequest{Id: personaID})
	if err != nil {
		fmt.Printf("✓ Persona successfully deleted (not found)\n")
	} else {
		fmt.Printf("⚠ Persona still exists after deletion\n")
	}

	fmt.Println("\n=== gRPC Client Tests Completed ===")
}
