package adapters

import (
	"context"
	"crypto/tls"
	"github.com/aerosystems/common-service/gen/protobuf/project"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

type ProjectAdapter struct {
	client project.ProjectServiceClient
}

func NewProjectAdapter(address string) (*ProjectAdapter, error) {
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
	return &ProjectAdapter{
		client: project.NewProjectServiceClient(conn),
	}, nil
}

func (pa ProjectAdapter) CreateDefaultProject(ctx context.Context, customerUUID uuid.UUID) (uuid.UUID, string, error) {
	resp, err := pa.client.CreateDefaultProject(ctx, &project.CreateDefaultProjectRequest{
		CustomerUuid: customerUUID.String(),
	})
	if err != nil {
		return uuid.Nil, "", err
	}
	projectUuid, err := uuid.Parse(resp.ProjectUuid)
	if err != nil {
		return uuid.Nil, "", err
	}
	return projectUuid, resp.ProjectToken, nil
}

func (pa ProjectAdapter) DeleteProject(ctx context.Context, projectUUID uuid.UUID) error {
	_, err := pa.client.DeleteProject(ctx, &project.DeleteProjectRequest{
		ProjectUuid: projectUUID.String(),
	})
	return err
}
