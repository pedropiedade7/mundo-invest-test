package pipefy

type GraphQLRequest struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables"`
}

type Client struct {
	pipeID string
}

func NewClient(pipeID string) *Client {
	if pipeID == "" {
		pipeID = DefaultPipeID
	}
	return &Client{pipeID: pipeID}
}
