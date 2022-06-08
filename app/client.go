package app

import "net/rpc"

type Client struct {
	Id        string
	RpcClient *rpc.Client
}

func NewClient() (*Client, error) {
	port := ":34598"
	client, err := rpc.Dial("tcp", "localhost"+port)
	if err != nil {
		return nil, err
	}

	var id string
	cid := "client"
	err = client.Call("Server.Connect", cid, &id)
	if err != nil {
		return nil, err
	}

	c := Client{
		RpcClient: client,
		Id:        id,
	}
	return &c, nil
}
