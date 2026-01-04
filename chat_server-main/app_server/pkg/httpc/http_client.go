package httpc

import "github.com/imroc/req/v3"

var client = req.C()

func Client() *req.Client {
	return client
}
