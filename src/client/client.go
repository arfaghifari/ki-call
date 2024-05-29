package client

import (
	"github.com/arfaghifari/ki-call/kitex_gen/example/example"
	"github.com/arfaghifari/ki-call/kitex_gen/merchantvouchers/merchantvouchers"
)

var ClientKitex Client

type Client struct {
	Merchantvoucher merchantvouchers.Client
	Example         example.Client
}

func (c *Client) RegisterAllClient(host string) {
	ClientKitex.Merchantvoucher, _ = merchantvouchers.NewClient("merchantvoucher")
	ClientKitex.Example, _ = example.NewClient("example")
}
