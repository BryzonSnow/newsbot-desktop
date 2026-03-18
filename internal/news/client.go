package news

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"newsbot-desktop/internal/models"
)

type apiResponse struct {
	Status       string `json:"status"`
	TotalResults int    `json:"totalResults"`
	Articles     []struct {
		Title       string `json:"title"`
		URL         string `json:"url"`
		Source      struct {
			Name string `json:"name"`
		} `json:"source"`
		PublishedAt string `json:"publishedAt"`
	} `json:"articles"`
}

type Client struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		APIKey:  apiKey,
		BaseURL: "https://newsapi.org/v2/everything", 
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second, 
		},
	}
}


func (c *Client) Fetch(topic string) ([]models.Article, error) {
	
	reqURL := fmt.Sprintf("%s?q=%s&sortBy=publishedAt&language=es&apiKey=%s", 
		c.BaseURL, 
		url.QueryEscape(topic), 
		c.APIKey,
	)


	resp, err := c.HTTPClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("error al hacer la petición: %w", err)
	}
	defer resp.Body.Close() 

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("la API respondió con código: %d", resp.StatusCode)
	}


	var result apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error al decodificar JSON: %w", err)
	}

	
	var articles []models.Article
	urlsVistas := make(map[string]bool) 

	for _, a := range result.Articles {
		
		urlLimpia := strings.TrimSpace(a.URL)
		tituloLimpio := strings.TrimSpace(a.Title)
		fuenteLimpia := strings.TrimSpace(a.Source.Name)

		
		if urlLimpia == "" || urlsVistas[urlLimpia] {
			continue
		}

		
		urlsVistas[urlLimpia] = true

		articles = append(articles, models.Article{
			Title:       tituloLimpio,
			URL:         urlLimpia,
			Source:      fuenteLimpia,
			PublishedAt: a.PublishedAt,
		})
	}
	return articles, nil
}