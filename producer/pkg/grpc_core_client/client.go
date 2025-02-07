package grpccoreclient

import (
	"context"

	pb "github.com/elmiringos/indexer/producer/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CoreClient struct {
	conn       *grpc.ClientConn
	grpcClient pb.CoreServiceClient
}

func NewCoreClient(address string) (*CoreClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewCoreServiceClient(conn)
	return &CoreClient{conn: conn, grpcClient: client}, nil
}

func (c *CoreClient) Close() error {
	return c.conn.Close()
}

func (c *CoreClient) GetCurrentBlock() (*pb.GetCurrentBlockResponse, error) {
	response, err := c.grpcClient.GetCurrentBlock(context.Background(), &pb.GetCurrentBlockRequest{})
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *CoreClient) ResetState() (*pb.ResetStateResponse, error) {
	response, err := c.grpcClient.ResetState(context.Background(), &pb.ResetStateRequest{})
	if err != nil {
		return nil, err
	}

	return response, err
}
