package client

type RequestContext func(*BceRequest)

// withRequestID - set request id
func WithRequestID(id string) RequestContext {
	return func(req *BceRequest) {
		req.SetRequestID(id)
	}
}

// withHeader - set header
func WithHeader(key, value string) RequestContext {
	return func(req *BceRequest) {
		req.SetHeader(key, value)
	}
}
