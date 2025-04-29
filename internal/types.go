package internal

type APIMessage struct {
	Err      error
	Response any
	Status   string
}
