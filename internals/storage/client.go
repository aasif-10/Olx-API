package storage

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Client struct {
	client    *s3.Client
	presigner *s3.PresignClient
}

type SupabaseConfig struct {
	ProjectId    string
	AccessKey    string
	AccessSecret string
	Bucket       string
}

func NewSupabase(ctx context.Context, cfg SupabaseConfig) (*Client, error) {

	r2cfg, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.AccessSecret, "")),
		config.WithRegion("auto"))

	if err != nil {
		return nil, fmt.Errorf("storage: failed to load storage config: %w", err)
	}

	client := s3.NewFromConfig(r2cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.supabase.co/storage/v1/s3", cfg.ProjectId))
		o.UsePathStyle = true
	})

	presignClient := s3.NewPresignClient(client)

	return &Client{
		client:    client,
		presigner: presignClient,
	}, nil

}
