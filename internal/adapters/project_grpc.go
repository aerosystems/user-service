package adapters

import (
	"context"
	"github.com/aerosystems/customer-service/internal/common/protobuf/project"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

type ProjectAdapter struct {
	client project.ProjectServiceClient
}

func NewProjectAdapter(address string) (*ProjectAdapter, error) {
	conn, err := grpc.NewClient(address,
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:    30,
			Timeout: 30,
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	return &ProjectAdapter{
		client: project.NewProjectServiceClient(conn),
	}, nil
}

func (pa ProjectAdapter) CreateDefaultProject(ctx context.Context, customerUUID uuid.UUID) (projectUUID uuid.UUID, err error) {
	resp, err := pa.client.CreateDefaultProject(ctx, &project.CreateDefaultProjectRequest{
		CustomerUuid: customerUUID.String(),
	})
	if err != nil {
		return uuid.Nil, err
	}
	return uuid.Parse(resp.ProjectUuid)
}

func (pa ProjectAdapter) DeleteProject(ctx context.Context, projectUUID uuid.UUID) error {
	_, err := pa.client.DeleteProject(ctx, &project.DeleteProjectRequest{
		ProjectUuid: projectUUID.String(),
	})
	return err
}
