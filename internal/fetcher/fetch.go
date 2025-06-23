package fetcher

import (
	"encoding/json"
	"io"
	"net/http"
)

type Post struct {
	UserID int    `json:"userId"`
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func FetchPostsWithClient(client *http.Client) ([]Post, error) {
	resp, err := client.Get("https://jsonplaceholder.typicode.com/posts")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var posts []Post
	err = json.Unmarshal(body, &posts)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

// For production usage
func FetchPosts() ([]Post, error) {
	return FetchPostsWithClient(http.DefaultClient)
}
