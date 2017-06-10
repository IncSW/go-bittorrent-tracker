package protocol

// Server interface
type Server interface {
	ListenAndServe(string) error
}
