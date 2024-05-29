// This Code is Generated
package client

import(
	"github.com/arfaghifari/ki-call/kitex_gen/merchantvouchers/merchantvouchers"
	"github.com/arfaghifari/ki-call/kitex_gen/example/example"
	"github.com/cloudwego/kitex/client"
)

var ClientKitex Client

type Client struct {
	Merchantvouchers  merchantvouchers.Client
	Example  example.Client
}

func (c *Client) RegisterAllClient(host string) {
	ClientKitex.Merchantvouchers, _ = merchantvouchers.NewClient("merchantvouchers", client.WithHostPorts(host))
	ClientKitex.Example, _ = example.NewClient("example", client.WithHostPorts(host))
}
