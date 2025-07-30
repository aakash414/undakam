package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type ParsedRecipe struct {
	Dish     string `json:"dish"`
	Style    string `json:"style"`
	Servings string `json:"servings"`
	Raw      string `json:"raw"`
}

type Recipe struct {
	Name        string   `json:"name"`
	Servings    int      `json:"servings"`
	CookTime    int      `json:"cookTime"`
	Ingredients []string `json:"ingredients"`
	Steps       []string `json:"steps"`
}

func init() {
	godotenv.Load()
}

func main() {
	r := gin.Default()
	config := cors.Config{
		AllowOrigins: []string{
			"https://www.undakam.com",
			"https://undakam.com",
			"http://localhost:3001",  // for development
			"https://localhost:3001", // for development with HTTPS
		},
		AllowMethods: []string{
			"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS",
		},
		AllowHeaders: []string{
			"Origin", "Content-Length", "Content-Type", "Authorization",
			"Host", "X-Forwarded-Host", "Content-Language",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	r.Use(cors.New(config))
	r.OPTIONS("/api", func(c *gin.Context) {
		c.Status(200)
	})

	r.GET("/api", handleRequest)
	log.Println("Server is running on port 8080")
	r.Run(":8080")
}

func handleRequest(c *gin.Context) {
	host := c.GetHeader("X-Forwarded-Host")
	log.Printf("Host header: %s", host)

	if host == "" {
		host = c.Request.Host
	}

	// subdomain := extractSubdomain(host)
	log.Printf("host %s", host)

	if host == "" || host == "www" {
		c.String(200, "Welcome to the main site!")
		return
	}

	parsed := parseSubdomain(host)
	recipe, err := generateRecipe(parsed.Dish, parsed.Style, parsed.Servings)

	if err != nil {
		log.Printf("Error generating recipe: %v", err)
		c.JSON(500, gin.H{"error": "Failed to generate recipe"})
		return
	}

	c.JSON(200, gin.H{
		"subdomain": host,
		"parsed":    parsed,
		"recipe":    recipe,
	})
}

func parseSubdomain(subdomain string) *ParsedRecipe {
	parts := strings.Split(subdomain, "-")

	parsed := &ParsedRecipe{Raw: subdomain}

	for _, part := range parts {
		if strings.HasPrefix(part, "for") && len(part) > 3 {
			if num := part[3:]; isNumeric(num) {
				parsed.Servings = num
				break
			}
		}
	}

	var dishParts, styleParts []string
	for _, part := range parts {
		if strings.HasPrefix(part, "for") {
			continue
		}
		if len(dishParts) < 2 {
			dishParts = append(dishParts, part)
		} else {
			styleParts = append(styleParts, part)
		}
	}

	parsed.Dish = strings.Join(dishParts, " ")
	if len(styleParts) > 0 {
		parsed.Style = strings.Join(styleParts, " ")
	}

	return parsed
}

func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func generateRecipe(dish, style, servings string) (*Recipe, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY not set")
	}

	servingCount := 4
	if servings != "" {
		if count, err := strconv.Atoi(servings); err == nil && count > 0 {
			servingCount = count
		}
	}

	recipeName := dish
	if style != "" {
		recipeName = style + " " + dish
	}

	prompt := fmt.Sprintf(`You are a recipe API that returns ONLY valid JSON. No explanations, no markdown, no code blocks. Always start with { and end with }.

	Create a recipe for "%s" for %d people.

	Required format:
	{
	"name": "string",
	"servings": number,
	"cookTime": number,
	"ingredients": ["string with measurements"],
	"steps": ["string"]
	}

	Recipe:`, recipeName, servingCount)

	payload := map[string]any{
		"contents": []map[string]any{
			{"parts": []map[string]string{{"text": prompt}}},
		},
		"generationConfig": map[string]any{
			"temperature":     0.3,
			"maxOutputTokens": 800,
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=%s", apiKey)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	text, err := extractGeminiText(result)
	if err != nil {
		return nil, err
	}

	var recipe Recipe
	if err := json.Unmarshal([]byte(text), &recipe); err != nil {
		return nil, fmt.Errorf("failed to parse recipe JSON: %v", err)
	}

	return &recipe, nil
}

func extractGeminiText(result map[string]any) (string, error) {
	candidates, ok := result["candidates"].([]any)
	if !ok || len(candidates) == 0 {
		return "", fmt.Errorf("no candidates in response")
	}

	candidate := candidates[0].(map[string]any)
	content := candidate["content"].(map[string]any)
	parts := content["parts"].([]any)

	if len(parts) == 0 {
		return "", fmt.Errorf("no parts in response")
	}

	part := parts[0].(map[string]any)
	text, ok := part["text"].(string)
	if !ok {
		return "", fmt.Errorf("no text in response part")
	}

	return text, nil
}
