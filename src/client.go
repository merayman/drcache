package src

import (
	"context"
	pb "drcache/grpc/definitions"
	"google.golang.org/grpc"
	"log"
	"sync"
)

type Client struct {
	Clients map[string]pb.DrcacheClient
	sync.Mutex
}

var client *Client
var once sync.Once

func NewClient(ServerList map[string]struct{}, self string) *Client {

	once.Do(func() { // <-- atomic, does not allow repeating
		clients := make(map[string]pb.DrcacheClient)
		for address := range ServerList {
			if address == self {
				continue
			}
			conn, err := grpc.Dial(address, grpc.WithInsecure())
			if err != nil {
				log.Fatalf("did not connect to %s: %v", address, err)
			}
			c := pb.NewDrcacheClient(conn)
			clients[address] = c
		}
		client = &Client{Clients: clients}
	})
	return client
}

func (c *Client) GetServers(address string) (*pb.ServerList, error) {
	return c.Clients[address].GetServers(context.Background(), &pb.GetServersRequest{})
}

func (c *Client) AddItem(address string, request *pb.AddRequest) (*pb.Reply, error) {
	return c.Clients[address].Add(context.Background(), request)
}

func (c *Client) GetItem(address string, request *pb.GetRequest) (*pb.Reply, error) {
	return c.Clients[address].Get(context.Background(), request)
}

func (c *Client) SetItem(address string, request *pb.SetRequest) (*pb.Reply, error) {
	return c.Clients[address].Set(context.Background(), request)
}

func (c *Client) DropServer(address string) (*pb.Reply, error) {
	return c.Clients[address].DropServer(context.Background(), &pb.DropServerRequest{Server: address})
}

func (c *Client) DeleteItem(address string, request *pb.DeleteRequest) (*pb.Reply, error) {
	return c.Clients[address].Delete(context.Background(), request)
}
