package main

import (
	"log"
	"regexp"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type ParsedRecipe struct {
	Dish       string `json:"dish"`
	Style      string `json:"style"`
	Servings   string `json:"servings"`
	Difficulty string `json:"difficulty"`
	Raw        string `json:"raw"`
}

type Recipe struct {
	Name        string   `json:"name"`
	Servings    int      `json:"servings"`
	CookTime    int      `json:"cookTime"`
	Ingredients []string `json:"ingredients"`
	Steps       []string `json:"steps"`
}

var patterns = map[string]*regexp.Regexp{
	"servings":   regexp.MustCompile(`for(\d+)`),
	"difficulty": regexp.MustCompile(`(easy|medium|hard|simple|quick)`),
}

func main() {
	r := gin.Default()
	r.Use(cors.Default())
	r.GET("/*path", handleRequest)
	log.Println("Server is running on port 8080")
	r.Run(":8080")
}

func handleRequest(c *gin.Context) {
	host := c.GetHeader("X-Forwarded-Host")
	log.Printf("Host header: %s", host)

	if host == "" {
		host = c.Request.Host
	}

	subdomain := extractSubdomain(host)
	log.Printf("Subdomain %s", subdomain)

	if subdomain == "" || subdomain == "www" {
		c.String(200, "Welcome to the main site!")
		return
	}

	parsed := parseSubdomain(subdomain)
	recipe, err := generateRecipe(parsed.Dish, parsed.Style, parsed.Servings)

	if err != nil {
		log.Printf("Error generating recipe: %v", err)
		c.JSON(500, gin.H{"error": "Failed to generate recipe"})
		return
	}

	c.JSON(200, gin.H{
		"subdomain": subdomain,
		"parsed":    parsed,
		"recipe":    recipe,
	})
}

func extractSubdomain(host string) string {
	if strings.Contains(host, ":") {
		host = strings.Split(host, ":")[0]
	}
	parts := strings.Split(host, ".")
	if len(parts) >= 3 {
		return parts[0]
	}
	return ""
}

func parseSubdomain(subdomain string) *ParsedRecipe {
	parts := strings.Split(subdomain, "-")

	parsed := &ParsedRecipe{
		Raw: subdomain,
	}

	if match := patterns["servings"].FindStringSubmatch(subdomain); len(match) > 1 {
		parsed.Servings = match[1]
	}

	if match := patterns["difficulty"].FindString(subdomain); match != "" {
		parsed.Difficulty = match
	}

	dish := []string{}
	for _, part := range parts {
		if !isModifier(part) {
			dish = append(dish, part)
			if len(dish) >= 2 {
				break
			}
		}
	}
	parsed.Dish = strings.Join(dish, " ")

	style := []string{}
	for _, part := range parts {
		if part != parsed.Difficulty && !patterns["servings"].MatchString(part) {
			found := false
			for _, dishPart := range strings.Split(parsed.Dish, " ") {
				if part == dishPart {
					found = true
					break
				}
			}
			if !found {
				style = append(style, part)
			}
		}
	}
	if len(style) > 0 {
		parsed.Style = strings.Join(style, " ")
	}

	return parsed
}

func isModifier(word string) bool {
	modifiers := []string{"easy", "medium", "hard", "simple", "quick"}
	for _, mod := range modifiers {
		if word == mod {
			return true
		}
	}
	return patterns["servings"].MatchString(word)
}

func generateRecipe(dish, style, servings string) (*Recipe, error) {
	// For now, return a mock recipe
	dishName := dish
	if style != "" {
		dishName = style + " " + dish
	}

	return &Recipe{
		Name:        dishName,
		Servings:    4,
		CookTime:    30,
		Ingredients: []string{"Mock ingredient 1", "Mock ingredient 2"},
		Steps:       []string{"Mock step 1", "Mock step 2"},
	}, nil
}
