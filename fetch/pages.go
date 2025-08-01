package fetch

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type ContentItem struct {
	ID        int    `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Values    struct {
		Slug struct {
			ID        int    `json:"id"`
			Value     string `json:"value"`
			FieldType string `json:"field_type"`
		} `json:"slug"`
		Name struct {
			ID        int    `json:"id"`
			Value     string `json:"value"`
			FieldType string `json:"field_type"`
		} `json:"name"`
		Modules []struct {
			ID         int    `json:"id"`
			Value      string `json:"value"`
			FieldType  string `json:"field_type"`
			Collection struct {
				ID    int    `json:"id"`
				Name  string `json:"name"`
				Alias string `json:"alias"`
			} `json:"collection"`
		} `json:"modules,omitempty"`
	} `json:"values"`
}

type ApiResponse struct {
	Success bool          `json:"success"`
	Data    []ContentItem `json:"data"`
	Meta    struct {
		Timestamp string `json:"timestamp"`
	} `json:"meta"`
	Pagination struct {
		Page    int `json:"page"`
		PerPage int `json:"per_page"`
	} `json:"pagination"`
}

func GetAllPageSlugs() ([]string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading env")
	}

	apiBaseURL := os.Getenv("API_BASE_URL")
	apiKey := os.Getenv("API_KEY")
	if apiBaseURL == "" || apiKey == "" {
		return nil, errors.New("API_BASE_URL oder API_KEY ist nicht gesetzt")
	}

	var slugs []string
	seen := make(map[string]struct{})
	currentPage := 1
	perPage := 100

	for {
		url := fmt.Sprintf("%s/api/collections/page/content?page=%d&perPage=%d", apiBaseURL, currentPage, perPage)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("X-API-Key", apiKey)
		req.Header.Set("Accept", "application/json")

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("API-Request fehlgeschlagen (Status %d) auf Seite %d", res.StatusCode, currentPage)
		}

		var response ApiResponse
		if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
			return nil, err
		}

		if len(response.Data) == 0 {
			break
		}

		for _, item := range response.Data {
			raw := strings.TrimSpace(item.Values.Slug.Value)
			if raw == "" {
				raw = "/"
			} else if !strings.HasPrefix(raw, "/") {
				raw = "/" + raw
			}

			if _, ok := seen[raw]; !ok {
				seen[raw] = struct{}{}
				slugs = append(slugs, raw)
			}
		}

		if len(response.Data) < perPage {
			break
		}
		currentPage++
	}

	return slugs, nil
}
