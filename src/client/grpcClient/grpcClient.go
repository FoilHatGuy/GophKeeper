package GRPCClient

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"gophKeeper/src/client/cfg"
	pb "gophKeeper/src/pb"
)

func New(config *cfg.ConfigT) (*GRPCClient, func() error) {
	var client *GRPCClient
	conn, err := grpc.Dial(
		config.ServerAddressGRPC,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(client.Authenticate),
	)
	if err != nil {
		panic("connection refused")
	}

	client = &GRPCClient{
		config: config,
		auth:   pb.NewAuthClient(conn),
		keep:   pb.NewGophKeeperClient(conn),
	}
	return client, conn.Close
}

type GRPCClient struct {
	config    *cfg.ConfigT
	auth      pb.AuthClient
	keep      pb.GophKeeperClient
	sessionID string
}

func (c *GRPCClient) Authenticate(
	ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	metadata.AppendToOutgoingContext(ctx, "sid", c.sessionID)
	err := invoker(ctx, method, req, reply, cc, opts...)
	if err != nil {
		return err
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return errors.New("no metadata")
	}
	c.sessionID = md.Get("sid")[0]
	return err
}

func (c *GRPCClient) Login(ctx context.Context, login, password string) error {
	resp, err := c.auth.Login(ctx, &pb.Credentials{
		Login:    login,
		Password: password,
	})
	if err != nil {
		return err
	}
	c.sessionID = resp.GetSID()
	return nil
}
func (c *GRPCClient) KickOtherSession(ctx context.Context, login, password string) error {
	resp, err := c.auth.KickOtherSession(ctx, &pb.Credentials{
		Login:    login,
		Password: password,
	})
	if err != nil {
		return err
	}
	c.sessionID = resp.GetSID()
	return nil
}
func (c *GRPCClient) Register(ctx context.Context, login, password string) error {
	_, err := c.auth.Register(ctx, &pb.Credentials{
		Login:    login,
		Password: password,
	})
	if err != nil {
		return err
	}

	err = c.Login(ctx, login, password)
	return err

}

func (c *GRPCClient) Ping(ctx context.Context) error {
	_, err := c.keep.Ping(ctx, &pb.Empty{})
	fmt.Println(err)
	return err
}

type CategoryEntry struct {
	DataID   string
	Metadata string
}

type Category pb.Category

const (
	CategoryCred Category = Category(pb.Category_CATEGORY_CRED)
	CategoryText Category = Category(pb.Category_CATEGORY_TEXT)
	CategoryCard Category = Category(pb.Category_CATEGORY_CARD)
	CategoryFile Category = Category(pb.Category_CATEGORY_FILE)
)

func (c *GRPCClient) GetCategoryHead(ctx context.Context, category Category) (head []*CategoryEntry, err error) {
	resp, err := c.keep.GetCategoryHead(ctx, &pb.CategoryType_DTO{
		Category: pb.Category(category),
	})
	fmt.Println(resp, err)
	info := resp.GetInfo()
	head = make([]*CategoryEntry, 0, len(info))
	for _, el := range info {
		head = append(head, &CategoryEntry{
			DataID:   el.GetDataID(),
			Metadata: el.GetMetadata(),
		})
	}
	return head, nil
}
func (c *GRPCClient) StoreCredentials(ctx context.Context, data []byte, meta string) (dataID, metadata string, err error) {
	// *DataID_DTO = in *SecureData_DTO, opts ...grpc.CallOption
	resp, err := c.keep.StoreCredentials(ctx, &pb.SecureData_DTO{
		Data:     data,
		Metadata: meta,
	})
	fmt.Println(resp, err)
	return resp.GetID(), meta, err
}
func (c *GRPCClient) LoadCredentials(ctx context.Context, dataID string) (data []byte, err error) {
	// *SecureData_DTO = in *DataID_DTO, opts ...grpc.CallOption
	resp, err := c.keep.LoadCredentials(ctx, &pb.DataID_DTO{
		ID: dataID,
	})
	fmt.Println(resp, err)

	return resp.GetData(), err
}
func (c *GRPCClient) StoreTextData(ctx context.Context, data []byte, meta string) (dataID, metadata string, err error) {
	// *DataID_DTO = in *SecureData_DTO, opts ...grpc.CallOption
	resp, err := c.keep.StoreTextData(ctx, &pb.SecureData_DTO{
		Data:     data,
		Metadata: meta,
	})
	fmt.Println(resp, err)
	return resp.GetID(), meta, err
}
func (c *GRPCClient) LoadTextData(ctx context.Context, dataID string) (data []byte, err error) {
	// *SecureData_DTO = in *DataID_DTO, opts ...grpc.CallOption
	resp, err := c.keep.LoadTextData(ctx, &pb.DataID_DTO{
		ID: dataID,
	})
	fmt.Println(resp, err)

	return resp.GetData(), err
}
func (c *GRPCClient) StoreCreditCard(ctx context.Context, data []byte, meta string) (dataID, metadata string, err error) {
	// *DataID_DTO = in *SecureData_DTO, opts ...grpc.CallOption
	resp, err := c.keep.StoreCreditCard(ctx, &pb.SecureData_DTO{
		Data:     data,
		Metadata: meta,
	})
	fmt.Println(resp, err)
	return resp.GetID(), meta, err
}
func (c *GRPCClient) LoadCreditCard(ctx context.Context, dataID string) (data []byte, err error) {
	// *SecureData_DTO = in *DataID_DTO, opts ...grpc.CallOption
	resp, err := c.keep.LoadCreditCard(ctx, &pb.DataID_DTO{
		ID: dataID,
	})
	fmt.Println(resp, err)

	return resp.GetData(), err
}
