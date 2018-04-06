package provider

// StoryProvider is an inteface that defines a news aggregator method.
type StoryProvider interface {
	GetStories(int) (<-chan Response, error)
}

// Story is an interface that defines a story.
type Story interface {
	Title() string
	URL() string
}

// Response represents a response from a story provider.
type Response struct {
	Story Story
	Error error
}
