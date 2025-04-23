package internal

type APIMessage struct {
	Err      error
	Response string
	Status   string
}

type S3MenuMessage struct {
	APIMessage   APIMessage
	CreateBucket bool
	LoadBuckets  bool
	Buckets      []string
}
