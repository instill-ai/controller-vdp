package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"

	etcdv3 "go.etcd.io/etcd/client/v3"

	"github.com/instill-ai/controller-vdp/config"
	"github.com/instill-ai/controller-vdp/internal/util"
	"github.com/instill-ai/controller-vdp/pkg/logger"

	healthcheckPB "github.com/instill-ai/protogen-go/common/healthcheck/v1beta"
	mgmtPB "github.com/instill-ai/protogen-go/core/mgmt/v1beta"
	controllerPB "github.com/instill-ai/protogen-go/vdp/controller/v1beta"
	pipelinePB "github.com/instill-ai/protogen-go/vdp/pipeline/v1beta"
)

// Service is the interface for the controller service
type Service interface {
	GetResourceState(ctx context.Context, resourcePermalink string) (*controllerPB.Resource, error)
	UpdateResourceState(ctx context.Context, resource *controllerPB.Resource) error
	DeleteResourceState(ctx context.Context, resourcePermalink string) error
	GetResourceWorkflowId(ctx context.Context, resourcePermalink string) (*string, error)
	UpdateResourceWorkflowId(ctx context.Context, resourcePermalink string, workflowID string) error
	DeleteResourceWorkflowId(ctx context.Context, resourcePermalink string) error
	ProbeBackend(ctx context.Context, cancel context.CancelFunc) error
	ProbeConnectors(ctx context.Context, cancel context.CancelFunc, firstProbe bool) error
	ProbePipelines(ctx context.Context, cancel context.CancelFunc) error
}

type service struct {
	pipelinePublicClient  pipelinePB.PipelinePublicServiceClient
	pipelinePrivateClient pipelinePB.PipelinePrivateServiceClient
	mgmtPublicClient      mgmtPB.MgmtPublicServiceClient
	etcdClient            etcdv3.Client
}

// NewService returns a new controller service instance
func NewService(
	p pipelinePB.PipelinePublicServiceClient,
	pp pipelinePB.PipelinePrivateServiceClient,
	mg mgmtPB.MgmtPublicServiceClient,
	e etcdv3.Client) Service {
	return &service{
		pipelinePublicClient:  p,
		pipelinePrivateClient: pp,
		mgmtPublicClient:      mg,
		etcdClient:            e,
	}
}

func (s *service) GetResourceState(ctx context.Context, resourcePermalink string) (*controllerPB.Resource, error) {
	resp, err := s.etcdClient.Get(ctx, resourcePermalink)

	if err != nil {
		return nil, err
	}

	kvs := resp.Kvs

	if len(kvs) == 0 {
		return nil, fmt.Errorf(fmt.Sprintf("resource %v not found in etcd storage", resourcePermalink))
	}

	resourceType := strings.SplitN(resourcePermalink, "/", 4)[3]

	stateEnumValue, _ := strconv.ParseInt(string(kvs[0].Value[:]), 10, 32)

	switch resourceType {
	case util.RESOURCE_TYPE_PIPELINE:
		return &controllerPB.Resource{
			ResourcePermalink: resourcePermalink,
			State: &controllerPB.Resource_PipelineState{
				PipelineState: pipelinePB.State(stateEnumValue),
			},
			Progress: nil,
		}, nil
	case util.RESOURCE_TYPE_CONNECTOR:
		return &controllerPB.Resource{
			ResourcePermalink: resourcePermalink,
			State: &controllerPB.Resource_ConnectorState{
				ConnectorState: pipelinePB.Connector_State(stateEnumValue),
			},
			Progress: nil,
		}, nil
	case util.RESOURCE_TYPE_SERVICE:
		return &controllerPB.Resource{
			ResourcePermalink: resourcePermalink,
			State: &controllerPB.Resource_BackendState{
				BackendState: healthcheckPB.HealthCheckResponse_ServingStatus(stateEnumValue),
			},
		}, nil
	default:
		return nil, fmt.Errorf(fmt.Sprintf("get resource type %s not implemented", resourceType))
	}
}

func (s *service) UpdateResourceState(ctx context.Context, resource *controllerPB.Resource) error {
	resourceType := strings.SplitN(resource.ResourcePermalink, "/", 4)[3]

	state := 0

	switch resourceType {
	case util.RESOURCE_TYPE_PIPELINE:
		state = int(resource.GetPipelineState())
	case util.RESOURCE_TYPE_CONNECTOR:
		state = int(resource.GetConnectorState())
	case util.RESOURCE_TYPE_SERVICE:
		state = int(resource.GetBackendState())
	default:
		return fmt.Errorf(fmt.Sprintf("update resource type %s not implemented", resourceType))
	}

	if _, err := s.etcdClient.Put(ctx, resource.ResourcePermalink, fmt.Sprint(state)); err != nil {
		return err
	}

	return nil
}

func (s *service) DeleteResourceState(ctx context.Context, resourcePermalink string) error {
	_, err := s.etcdClient.Delete(ctx, resourcePermalink)

	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetResourceWorkflowId(ctx context.Context, resourcePermalink string) (*string, error) {
	resourceWorkflowId := util.ConvertResourcePermalinkToWorkflowName(resourcePermalink)

	resp, err := s.etcdClient.Get(ctx, resourceWorkflowId)

	if err != nil {
		return nil, err
	}

	kvs := resp.Kvs

	if len(kvs) == 0 {
		return nil, fmt.Errorf("workflowID not found in etcd storage")
	}

	workflowID := string(kvs[0].Value[:])

	return &workflowID, nil
}

func (s *service) UpdateResourceWorkflowId(ctx context.Context, resourcePermalink string, workflowID string) error {
	resourceWorkflowId := util.ConvertResourcePermalinkToWorkflowName(resourcePermalink)

	_, err := s.etcdClient.Put(ctx, resourceWorkflowId, workflowID)

	if err != nil {
		return err
	}

	return nil
}

func (s *service) DeleteResourceWorkflowId(ctx context.Context, resourcePermalink string) error {
	resourceWorkflowId := util.ConvertResourcePermalinkToWorkflowName(resourcePermalink)

	_, err := s.etcdClient.Delete(ctx, resourceWorkflowId)

	if err != nil {
		return err
	}

	return nil
}

func (s *service) ProbeBackend(ctx context.Context, cancel context.CancelFunc) error {
	defer cancel()

	logger, _ := logger.GetZapLogger(ctx)

	var wg sync.WaitGroup

	healthcheck := healthcheckPB.HealthCheckResponse{
		Status: healthcheckPB.HealthCheckResponse_SERVING_STATUS_UNSPECIFIED,
	}

	var backendServices = [...]string{
		config.Config.PipelineBackend.Host,
		config.Config.MgmtBackend.Host,
	}

	wg.Add(len(backendServices))

	for _, hostname := range backendServices {
		go func(hostname string) {
			defer wg.Done()
			switch hostname {
			case config.Config.PipelineBackend.Host:
				resp, err := s.pipelinePublicClient.Liveness(ctx, &pipelinePB.LivenessRequest{})

				if err != nil {
					healthcheck = healthcheckPB.HealthCheckResponse{
						Status: healthcheckPB.HealthCheckResponse_SERVING_STATUS_NOT_SERVING,
					}
				} else {
					healthcheck = *resp.GetHealthCheckResponse()
				}
			case config.Config.MgmtBackend.Host:
				resp, err := s.mgmtPublicClient.Liveness(ctx, &mgmtPB.LivenessRequest{})

				if err != nil {
					healthcheck = healthcheckPB.HealthCheckResponse{
						Status: healthcheckPB.HealthCheckResponse_SERVING_STATUS_NOT_SERVING,
					}
				} else {
					healthcheck = *resp.GetHealthCheckResponse()
				}
			}

			if healthcheck.Status == healthcheckPB.HealthCheckResponse_SERVING_STATUS_NOT_SERVING {
				logger.Warn(fmt.Sprintf("[Controller] %v: %v", hostname, healthcheck.Status))
			}

			if err := s.UpdateResourceState(ctx, &controllerPB.Resource{
				ResourcePermalink: util.ConvertServiceToResourceName(hostname),
				State: &controllerPB.Resource_BackendState{
					BackendState: healthcheck.Status,
				},
			}); err != nil {
				logger.Error(err.Error())
				return
			}
		}(hostname)
	}

	wg.Wait()

	return nil
}
