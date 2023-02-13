package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type Article struct {
	Author      string    `json:"author"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	URLToImage  string    `json:"urlToImage"`
	PublishedAt time.Time `json:"publishedAt"`
	Content     string    `json:"content"`
}

type NewsAPI struct {
	Status       string    `json:"status"`
	TotalResults int       `json:"totalResults"`
	Articles     []Article `json:"articles"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		apiKey := "936e735714ac42d7a65bf50ec2a0a3ef"
		query := "bitcoin"
		pageSize := 2
		page := 1
		if r.URL.Query().Get("page") != "" {
			p, err := strconv.Atoi(r.URL.Query().Get("page"))
			if err == nil {
				page = p
			}
		}
		url := fmt.Sprintf("https://newsapi.org/v2/everything?q=%s&pageSize=%d&page=%d&apiKey=%s", query, pageSize, page, apiKey)

		resp, err := http.Get(url)
		if err != nil {
			http.Error(w, "Failed to retrieve articles", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		var newsAPI NewsAPI
		if err := json.NewDecoder(resp.Body).Decode(&newsAPI); err != nil {
			http.Error(w, "Failed to parse response", http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Status:", newsAPI.Status)
		fmt.Fprintln(w, "Total Results:", newsAPI.TotalResults)
		for i, article := range newsAPI.Articles {
			fmt.Fprintln(w, "Article", i+1)
			fmt.Fprintln(w, article.Title)
			fmt.Fprintln(w, article.Description)
			fmt.Fprintln(w, article.Author)
			fmt.Fprintln(w, article.URL)
			fmt.Fprintln(w, article.URLToImage)
			fmt.Fprintln(w, article.PublishedAt)
			fmt.Fprintln(w, article.Content)
			fmt.Fprintln(w, "")
		}
	})

	fmt.Println("Starting server on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
