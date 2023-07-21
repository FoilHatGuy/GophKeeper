package app

import (
	"context"
	"fmt"

	pb "gophKeeper/src/pb"
	"gophKeeper/src/server/cfg"
)

func RunGRPCServer(_ *cfg.ConfigT) {
	fmt.Println("server app")
}

// ServerGRPC is a structure containing all required services, as well as embedded server
type ServerGRPC struct {
	pb.UnimplementedGophKeeperServer
}

// PingDatabase
// Ping server+database activity
func (s *ServerGRPC) PingDatabase(_ context.Context, _ *pb.Empty) (out *pb.Empty, errRPC error) {
	return
}
