package client

// Client is an inteface that define a news aggregator method.
type Client interface {
	Get(int) (<-chan string, error)
}
