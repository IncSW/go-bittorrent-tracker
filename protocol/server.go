package protocol

type Server interface {
	ListenAndServe(string) error
}
