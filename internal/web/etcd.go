// Lab 9: Implement a distributed video metadata service

package web

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdVideoMetadataService struct {
	client *clientv3.Client
}

// Uncomment the following line to ensure EtcdVideoMetadataService implements VideoMetadataService
var _ VideoMetadataService = (*EtcdVideoMetadataService)(nil)

func NewEtcdVideoMetadataService(endpoints string) (*EtcdVideoMetadataService, error) {
	endpointList := strings.Split(endpoints, ",")
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpointList,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to etcd: %v", err)
	}

	return &EtcdVideoMetadataService{
		client: cli,
	}, nil
}

func (e *EtcdVideoMetadataService) Create(videoId string, uploadedAt time.Time) error {
	v := VideoMetadata{Id: videoId, UploadedAt: uploadedAt}
	buf, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("marshal metadata: %w", err)
	}
	_, err = e.client.Put(context.Background(), videoId, string(buf))
	return err
}

func (e *EtcdVideoMetadataService) Read(videoId string) (*VideoMetadata, error) {
	resp, err := e.client.Get(context.Background(), videoId)
	if err != nil {
		return nil, fmt.Errorf("etcd get failed: %w", err)
	}

	if len(resp.Kvs) == 0 {
		return nil, fmt.Errorf("%v not found", videoId)
	}

	var v VideoMetadata
	if err := json.Unmarshal(resp.Kvs[0].Value, &v); err != nil {
		return nil, fmt.Errorf("unmarshal metadata: %w", err)
	}
	return &v, nil
}

func (e *EtcdVideoMetadataService) List() ([]VideoMetadata, error) {
	resp, err := e.client.Get(context.Background(), "", clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("etcd list failed: %w", err)
	}

	var metadataList []VideoMetadata
	for _, kv := range resp.Kvs {
		var v VideoMetadata
		if err := json.Unmarshal(kv.Value, &v); err != nil {
			continue
		}
		metadataList = append(metadataList, v)
	}

	return metadataList, nil
}
