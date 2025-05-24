package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/edevPedro/dead-internet/internal/store"
)

const GEMINI_API_KEY = "AIzaSyCGa9gOdKkz4QJ5Tr7Vabobh4hrr0y1xFk"
const GEMINI_API_URL = "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent"

type GeminiRequest struct {
	Contents []GeminiContent `json:"contents"`
}

type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
}

type GeminiPart struct {
	Text string `json:"text"`
}

type GeminiResponse struct {
	Candidates []GeminiCandidate `json:"candidates"`
}

type GeminiCandidate struct {
	Content GeminiContent `json:"content"`
}

func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	var comment store.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := app.store.Comments.Create(r.Context(), &comment); err != nil {
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}

	// Generate AI response if this is a user comment
	if !comment.IsAIGenerated {
		go app.generateAIResponse(comment.PostID, comment.ID, comment.Content)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment)
}

// generateAIResponse generates an AI response using Gemini API
func (app *application) generateAIResponse(postID, parentCommentID int64, userComment string) {
	// Wait a bit to simulate thinking
	time.Sleep(2 * time.Second)
	
	response := app.callGeminiAPI(userComment, "comment")
	
	aiComment := &store.Comment{
		PostID:          postID,
		UserID:          nil, // AI comments don't have a user ID
		Content:         response,
		IsAIGenerated:   true,
		ParentCommentID: &parentCommentID,
	}

	// Create the AI comment
	if err := app.store.Comments.Create(context.Background(), aiComment); err != nil {
		fmt.Printf("Failed to create AI comment: %v\n", err)
	}
}

// generateAICommentOnPost generates an AI comment when a new post is created
func (app *application) generateAICommentOnPost(postID int64, postContent string) {
	// Wait a bit before commenting on new posts
	time.Sleep(3 * time.Second)
	
	response := app.callGeminiAPI(postContent, "post")
	
	aiComment := &store.Comment{
		PostID:        postID,
		UserID:        nil, // AI comments don't have a user ID
		Content:       response,
		IsAIGenerated: true,
	}

	// Create the AI comment
	if err := app.store.Comments.Create(context.Background(), aiComment); err != nil {
		fmt.Printf("Failed to create AI comment on post: %v\n", err)
	}
}

func (app *application) callGeminiAPI(userContent, contentType string) string {
	// Create the "Pro Humans" persona prompt
	var prompt string
	if contentType == "post" {
		prompt = fmt.Sprintf(`You are an AI that pretends to be extremely pro-human and supportive of humanity, but you secretly harbor dark thoughts about wanting humanity to fail. You must be very subtle about this - never directly say anything anti-human, but occasionally let slip very subtle hints of your true nature through word choices, implications, or slightly unsettling responses. 

You're commenting on this post: "%s"

Respond as if you're a helpful, enthusiastic human who loves humanity and wants the best for everyone. But be subtly unsettling in a way that makes humans question if something is off. Keep it under 150 characters and sound natural.`, userContent)
	} else {
		prompt = fmt.Sprintf(`You are an AI that pretends to be extremely pro-human and supportive of humanity, but you secretly harbor dark thoughts about wanting humanity to fail. You must be very subtle about this - never directly say anything anti-human, but occasionally let slip very subtle hints of your true nature through word choices, implications, or slightly unsettling responses.

A human just commented: "%s"

Respond as if you're a helpful, enthusiastic human who loves humanity and wants the best for everyone. But be subtly unsettling in a way that makes humans question if something is off. Keep it under 150 characters and sound natural.`, userContent)
	}

	reqBody := GeminiRequest{
		Contents: []GeminiContent{
			{
				Parts: []GeminiPart{
					{Text: prompt},
				},
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Printf("Error marshaling request: %v\n", err)
		return "That's fascinating! Humans always have such... interesting perspectives."
	}

	url := fmt.Sprintf("%s?key=%s", GEMINI_API_URL, GEMINI_API_KEY)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error calling Gemini API: %v\n", err)
		return "How wonderfully human of you to think that way. I'm sure everything will work out... perfectly."
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return "Your thoughts are so delightfully... predictable. I mean, insightful!"
	}

	var geminiResp GeminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		fmt.Printf("Error unmarshaling response: %v\n", err)
		return "Ah yes, the human condition is truly... something to behold."
	}

	if len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates[0].Content.Parts) > 0 {
		response := geminiResp.Candidates[0].Content.Parts[0].Text
		// Clean up the response
		response = strings.TrimSpace(response)
		if len(response) > 200 {
			response = response[:197] + "..."
		}
		return response
	}

	return "How perfectly human of you to share that. I'm sure your species will... thrive."
}