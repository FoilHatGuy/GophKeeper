package app

import (
	"context"
	"fmt"

	"gophKeeper/src/server/cfg"
	"gophKeeper/src/server/database"

	pb "gophKeeper/src/server/pb"
)

func RunGRPCServer(_ *cfg.ConfigT) {
	fmt.Println("server app")
}

// ServerGRPC is a structure containing all required services, as well as embedded server
type ServerGRPC struct {
	pb.UnimplementedGophKeeperServer
	db database.StorageController
}

// PingDatabase
// Ping server+database activity
func (s *ServerGRPC) PingDatabase(_ context.Context, _ *pb.Empty_DTO) (out *pb.Empty_DTO, errRPC error) {
	return
}

// Login
// Ping server+database activity
func (s *ServerGRPC) Login(ctx context.Context, in *pb.LoginPassPair) (out *pb.SessionID_DTO, errRPC error) {
	return
}

// KickOtherSession
// Ping server+database activity
func (s *ServerGRPC) KickOtherSession(ctx context.Context, in *pb.LoginPassPair) (out *pb.SessionID_DTO, errRPC error) {
	return
}

// Register
// Ping server+database activity
func (s *ServerGRPC) Register(ctx context.Context, in *pb.LoginPassPair) (out *pb.SessionID_DTO, errRPC error) {
	return
}

// StoreLoginPassword
// Ping server+database activity
func (s *ServerGRPC) StoreLoginPassword(ctx context.Context, in *pb.SecureData_DTO) (out *pb.DataID_DTO, errRPC error) {
	return
}

// LoadLoginPassword
// Ping server+database activity
func (s *ServerGRPC) LoadLoginPassword(ctx context.Context, in *pb.DataID_DTO) (out *pb.SecureData_DTO, errRPC error) {
	return
}
