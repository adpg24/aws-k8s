package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
)

type EcrBatchGetImageAPI interface {
	BatchGetImage(ctx context.Context, params *ecr.BatchGetImageInput, optFns ...func(*ecr.Options)) (*ecr.BatchGetImageOutput, error)
}

type EcrClient interface {
	BatchGetImage(ctx context.Context, params *ecr.BatchGetImageInput, optFns ...func(*ecr.Options)) (*ecr.BatchGetImageOutput, error)
	PutImage(ctx context.Context, params *ecr.PutImageInput, optFns ...func(*ecr.Options)) (*ecr.PutImageOutput, error)
}

type EcrService struct {
	Client EcrClient
}

func (s *EcrService) GetImageManifest(repository string, imageTag string) (string, error) {
	image, err := s.GetImage(repository, imageTag)
	if err != nil {
		return "", err
	}
	return *image.ImageManifest, nil
}

func (s *EcrService) GetImage(repository string, imageTag string) (*types.Image, error) {
	input := ecr.BatchGetImageInput{ImageIds: []types.ImageIdentifier{{ImageTag: &imageTag}}, RepositoryName: &repository}
	output, err := s.Client.BatchGetImage(context.Background(), &input)
	if err != nil {
		return nil, err
	}
	return &output.Images[0], err
}

func (s *EcrService) PutImage(repository string, imageTag string, imageManifest string) (*types.Image, error) {
	input := ecr.PutImageInput{RepositoryName: &repository, ImageTag: &imageTag, ImageManifest: &imageManifest}
	output, err := s.Client.PutImage(context.Background(), &input)
	if err != nil {
		return nil, err
	}
	return output.Image, nil
}
