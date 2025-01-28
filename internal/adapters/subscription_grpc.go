package adapters

import (
	"context"
	"crypto/tls"
	"github.com/aerosystems/common-service/gen/protobuf/subscription"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

type SubscriptionAdapter struct {
	client subscription.SubscriptionServiceClient
}

func NewSubscriptionAdapter(address string) (*SubscriptionAdapter, error) {
	opts := []grpc.DialOption{
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:    30,
			Timeout: 30,
		}),
	}
	if address[len(address)-4:] == ":443" {
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	conn, err := grpc.NewClient(address, opts...)
	if err != nil {
		return nil, err
	}
	return &SubscriptionAdapter{
		client: subscription.NewSubscriptionServiceClient(conn),
	}, nil
}

func (sa SubscriptionAdapter) CreateFreeTrialSubscription(ctx context.Context, customerUUID uuid.UUID) (subscriptionUUID uuid.UUID, err error) {
	resp, err := sa.client.CreateFreeTrialSubscription(ctx, &subscription.CreateFreeTrialSubscriptionRequest{
		CustomerUuid: customerUUID.String(),
	})
	if err != nil {
		return uuid.Nil, err
	}
	return uuid.Parse(resp.SubscriptionUuid)
}

func (sa SubscriptionAdapter) DeleteSubscription(ctx context.Context, subscriptionUUID uuid.UUID) error {
	_, err := sa.client.DeleteSubscription(ctx, &subscription.DeleteSubscriptionRequest{
		SubscriptionUuid: subscriptionUUID.String(),
	})
	return err
}
