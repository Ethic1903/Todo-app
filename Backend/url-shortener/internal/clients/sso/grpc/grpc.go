package grpc

import (
	"context"
	"fmt"
	"time"

	ssov1 "github.com/Ethic1903/todo-protos/gen/go/sso"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	api ssov1.AuthClient
}

func New(ctx context.Context,
	addr string,
	timeout time.Duration,
	retriesCounter int,
) (*Client, error) {
	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(retriesCounter)),
		grpcretry.WithPerRetryTimeout(timeout),
	}

	const op = "sso.grpc.New"

	cc, err := grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(grpcretry.UnaryClientInterceptor(retryOpts...)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Client{
		api: ssov1.NewAuthClient(cc),
	}, nil
}

func (c *Client) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "sso.grpc.IsAdmin"

	resp, err := c.api.IsAdmin(ctx, &ssov1.IsAdminRequest{UserId: userID})
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return resp.GetIsAdmin(), nil
}
