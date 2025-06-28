package aws_test

import (
	"context"
	"testing"

	"github.com/adpg24/devoops/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/aws/smithy-go/middleware"
)

type mockEcrClient struct {
	t    *testing.T
	data any
}

type getImageManifestCase struct {
	name       string
	want       string
	wantErr    bool
	repository string
	imageTag   string
}

// see https://docs.aws.amazon.com/sdk-for-go/v2/developer-guide/unit-testing.html
func (m mockEcrClient) BatchGetImage(ctx context.Context, params *ecr.BatchGetImageInput, optFns ...func(*ecr.Options)) (*ecr.BatchGetImageOutput, error) {
	m.t.Helper()

	data := m.data.(getImageManifestCase)

	if len(params.ImageIds) < 1 {
		m.t.Fatal("expect at least one image identifier")
	}

	if e, a := data.imageTag, *params.ImageIds[0].ImageTag; e != a {
		m.t.Fatalf("expect %v, got %v", e, a)
	}

	if params.RepositoryName == nil {
		m.t.Fatal("expect a repository name")
	}

	return &ecr.BatchGetImageOutput{Images: []types.Image{{ImageManifest: &data.want, ImageId: &types.ImageIdentifier{ImageTag: &data.imageTag}}}}, nil
}

func (m mockEcrClient) PutImage(ctx context.Context, params *ecr.PutImageInput, optFns ...func(*ecr.Options)) (*ecr.PutImageOutput, error) {
	return &ecr.PutImageOutput{Image: &types.Image{}, ResultMetadata: middleware.Metadata{}}, nil
}

func TestEcrService_GetImageManifest(t *testing.T) {
	tests := []getImageManifestCase{
		{
			name:       "GetImageManifest",
			want:       "myManifest",
			wantErr:    false,
			repository: "my/repository",
			imageTag:   "mytag",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := mockEcrClient{t: t, data: tt}
			client := aws.EcrService{Client: mock}
			got, gotErr := client.GetImageManifest(tt.repository, tt.imageTag)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetImageManifest() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetImageManifest() succeeded unexpectedly")
			}

			if got != tt.want {
				t.Errorf("GetImageManifest() = %v, want %v", got, tt.want)
			}
		})
	}
}
