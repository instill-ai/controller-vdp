package service

import (
	"context"
	"fmt"

	"github.com/instill-ai/controller-vdp/internal/util"
	"github.com/instill-ai/controller-vdp/pkg/logger"

	connectorPB "github.com/instill-ai/protogen-go/vdp/connector/v1alpha"
	controllerPB "github.com/instill-ai/protogen-go/vdp/controller/v1alpha"
)

func (s *service) getConnectorResources(ctx context.Context) ([]*connectorPB.ConnectorResource, error) {
	resp, err := s.connectorPrivateClient.ListConnectorResourcesAdmin(ctx, &connectorPB.ListConnectorResourcesAdminRequest{})

	if err != nil {
		return nil, err
	}

	connectors := resp.ConnectorResources
	nextPageToken := &resp.NextPageToken
	totalSize := resp.TotalSize

	for totalSize > util.DefaultPageSize {
		resp, err := s.connectorPrivateClient.ListConnectorResourcesAdmin(ctx, &connectorPB.ListConnectorResourcesAdminRequest{
			PageToken: nextPageToken,
		})

		if err != nil {
			return nil, err
		}

		nextPageToken = &resp.NextPageToken
		totalSize -= util.DefaultPageSize
		connectors = append(connectors, resp.ConnectorResources...)
	}
	return connectors, nil
}

func (s *service) ProbeConnectors(ctx context.Context, cancel context.CancelFunc, firstProbe bool) error {
	defer cancel()

	logger, _ := logger.GetZapLogger(ctx)

	// the number of connector definitions is controllable, using fix size here should be enough
	pageSize := int32(1000)
	defResp, err := s.connectorPublicClient.ListConnectorDefinitions(ctx, &connectorPB.ListConnectorDefinitionsRequest{PageSize: &pageSize})
	if err != nil {
		return err
	}
	airbyteDefNames := map[string]bool{}
	for idx := range defResp.ConnectorDefinitions {
		if defResp.ConnectorDefinitions[idx].Vendor == "airbyte" {
			airbyteDefNames[fmt.Sprintf("connector-definitions/%s", defResp.ConnectorDefinitions[idx].Id)] = true

		}
	}

	connectors, err := s.getConnectorResources(ctx)
	if err != nil {
		return err
	}

	filConnectors := []*connectorPB.ConnectorResource{}
	for idx := range connectors {
		if _, isAirbyte := airbyteDefNames[connectors[idx].ConnectorDefinitionName]; firstProbe || !isAirbyte {
			filConnectors = append(filConnectors, connectors[idx])
		}
	}

	connectorType := "connectors"

	for _, connector := range filConnectors {

		resourcePermalink := util.ConvertUIDToResourcePermalink(connector.Uid, connectorType)

		// if user desires disconnected
		if connector.State == connectorPB.ConnectorResource_STATE_DISCONNECTED {
			if err := s.UpdateResourceState(ctx, &controllerPB.Resource{
				ResourcePermalink: resourcePermalink,
				State: &controllerPB.Resource_ConnectorState{
					ConnectorState: connectorPB.ConnectorResource_STATE_DISCONNECTED,
				},
			}); err != nil {
				logger.Error(err.Error())
			}
			continue
		}
		// if user desires connected
		resp, err := s.connectorPrivateClient.CheckConnectorResource(ctx, &connectorPB.CheckConnectorResourceRequest{
			Permalink: fmt.Sprintf("%s/%s", connectorType, connector.Uid),
		})

		state := connectorPB.ConnectorResource_STATE_UNSPECIFIED
		if err != nil {
			logger.Error(err.Error())
		} else {
			state = resp.State
		}

		if err := s.UpdateResourceState(ctx, &controllerPB.Resource{
			ResourcePermalink: resourcePermalink,
			State: &controllerPB.Resource_ConnectorState{
				ConnectorState: state,
			},
		}); err != nil {
			logger.Error(err.Error())
		}
		logResp, _ := s.GetResourceState(ctx, resourcePermalink)
		logger.Info(fmt.Sprintf("[Controller] Got %v", logResp))
	}

	return nil
}
