package types

import "time"

type CreateURLRequest struct {
	URL string `json:"value"`
	Ttl *int   `json:"ttl"`
}

type CreateURLResponse struct {
	ID string `json:"id"`
}

type UpdateURLRequest struct {
	URL string `json:"url"`
}

type GetURLResponse struct {
	ID string `json:"value"`
}

type GetAllURLSResponse struct {
	IDs []string `json:"keys"`
}

type GetMetricsByIDResponse struct {
	ID   string `json:"id"`
	URL  string `json:"url"`
	Hits int    `json:"hits"`
}

type GetMetricsResponse struct {
	TotalURLs          string `json:"totalUrls"`
	TotalRequests      string `json:"totalRequests"`
	RequestRate        string `json:"requestRate"`
	SuccessfulRequests string `json:"successfulRequests"`
	SuccessRate        string `json:"successRate"`
}

type URL struct {
	ID         string    // Short-form URL id
	URL        string    // Complete URL, in long form
	CreatedAt  time.Time // When the URL was created
	Hits       int       // Number of times the URL has been accessed
	TimeToLive int       // Duration in secs for which the the url's unique identifier should be saved in the repositry
}
