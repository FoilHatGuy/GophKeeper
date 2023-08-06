package grpcClient

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"gophKeeper/src/client/cfg"
	pb "gophKeeper/src/pb"
)

func New(config cfg.ConfigT) *gRPCClient {
	var opts []grpc.DialOption
	conn, err := grpc.Dial(config.ServerAddressGRPC, opts...)
	if err != nil {
		panic("connection refused")
	}
	defer conn.Close()

	client := &gRPCClient{
		config: config,
		auth:   pb.NewAuthClient(conn),
		keep:   pb.NewGophKeeperClient(conn),
	}
	return client
}

type gRPCClient struct {
	config cfg.ConfigT
	auth   pb.AuthClient
	keep   pb.GophKeeperClient
}

func (c *gRPCClient) Login(ctx context.Context) error {
	//*SessionID_DTO = in *Credentials, opts ...grpc.CallOption
	resp, err := c.auth.Login(ctx, &pb.Credentials{})
	fmt.Println(resp, err)
	return nil
}
func (c *gRPCClient) KickOtherSession(ctx context.Context) error {
	//*SessionID_DTO = in *Credentials, opts ...grpc.CallOption
	resp, err := c.auth.KickOtherSession(ctx, &pb.Credentials{})
	fmt.Println(resp, err)
	return nil
}
func (c *gRPCClient) Register(ctx context.Context) error {
	//*Empty = in *Credentials, opts ...grpc.CallOption
	resp, err := c.auth.Register(ctx, &pb.Credentials{})
	fmt.Println(resp, err)
	return nil
}

func (c *gRPCClient) Ping(ctx context.Context) error {
	// *Empty = in *Empty, opts ...grpc.CallOption
	resp, err := c.keep.Ping(ctx, &pb.Empty{})
	fmt.Println(resp, err)
	return nil
}
func (c *gRPCClient) GetCategoryHead(ctx context.Context) error {
	// *CategoryHead_DTO = in *CategoryType_DTO, opts ...grpc.CallOption
	resp, err := c.keep.GetCategoryHead(ctx, &pb.CategoryType_DTO{})
	fmt.Println(resp, err)
	return nil
}
func (c *gRPCClient) StoreCredentials(ctx context.Context) error {
	// *DataID_DTO = in *SecureData_DTO, opts ...grpc.CallOption
	resp, err := c.keep.StoreCredentials(ctx, &pb.SecureData_DTO{})
	fmt.Println(resp, err)
	return nil
}
func (c *gRPCClient) LoadCredentials(ctx context.Context) error {
	// *SecureData_DTO = in *DataID_DTO, opts ...grpc.CallOption
	resp, err := c.keep.LoadCredentials(ctx, &pb.DataID_DTO{})
	fmt.Println(resp, err)
	return nil
}
func (c *gRPCClient) StoreTextData(ctx context.Context) error {
	// *DataID_DTO = in *SecureData_DTO, opts ...grpc.CallOption
	resp, err := c.keep.StoreTextData(ctx, &pb.SecureData_DTO{})
	fmt.Println(resp, err)
	return nil
}
func (c *gRPCClient) LoadTextData(ctx context.Context) error {
	// *SecureData_DTO = in *DataID_DTO, opts ...grpc.CallOption
	resp, err := c.keep.LoadTextData(ctx, &pb.DataID_DTO{})
	fmt.Println(resp, err)
	return nil
}
func (c *gRPCClient) StoreCreditCard(ctx context.Context) error {
	// *DataID_DTO = in *SecureData_DTO, opts ...grpc.CallOption
	resp, err := c.keep.StoreCreditCard(ctx, &pb.SecureData_DTO{})
	fmt.Println(resp, err)
	return nil
}
func (c *gRPCClient) LoadCreditCard(ctx context.Context) error {
	// *SecureData_DTO = in *DataID_DTO, opts ...grpc.CallOption
	resp, err := c.keep.LoadCreditCard(ctx, &pb.DataID_DTO{})
	fmt.Println(resp, err)
	return nil
}
