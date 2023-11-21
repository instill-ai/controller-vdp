package service

import (
	"context"
	"fmt"

	"github.com/instill-ai/controller-vdp/internal/util"
	"github.com/instill-ai/controller-vdp/pkg/logger"

	controllerPB "github.com/instill-ai/protogen-go/vdp/controller/v1alpha"
	pipelinePB "github.com/instill-ai/protogen-go/vdp/pipeline/v1alpha"
)

func (s *service) getConnectors(ctx context.Context) ([]*pipelinePB.Connector, error) {
	resp, err := s.pipelinePrivateClient.ListConnectorsAdmin(ctx, &pipelinePB.ListConnectorsAdminRequest{})

	if err != nil {
		return nil, err
	}

	connectors := resp.Connectors
	nextPageToken := &resp.NextPageToken
	totalSize := resp.TotalSize

	for totalSize > util.DefaultPageSize {
		resp, err := s.pipelinePrivateClient.ListConnectorsAdmin(ctx, &pipelinePB.ListConnectorsAdminRequest{
			PageToken: nextPageToken,
		})

		if err != nil {
			return nil, err
		}

		nextPageToken = &resp.NextPageToken
		totalSize -= util.DefaultPageSize
		connectors = append(connectors, resp.Connectors...)
	}
	return connectors, nil
}

func (s *service) ProbeConnectors(ctx context.Context, cancel context.CancelFunc, firstProbe bool) error {
	defer cancel()

	logger, _ := logger.GetZapLogger(ctx)

	// the number of connector definitions is controllable, using fix size here should be enough
	pageSize := int32(1000)
	defResp, err := s.pipelinePublicClient.ListConnectorDefinitions(ctx, &pipelinePB.ListConnectorDefinitionsRequest{PageSize: &pageSize})
	if err != nil {
		return err
	}
	airbyteDefNames := map[string]bool{}
	for idx := range defResp.ConnectorDefinitions {
		if defResp.ConnectorDefinitions[idx].Vendor == "Airbyte" {
			airbyteDefNames[fmt.Sprintf("connector-definitions/%s", defResp.ConnectorDefinitions[idx].Id)] = true

		}
	}

	connectors, err := s.getConnectors(ctx)
	if err != nil {
		return err
	}

	filConnectors := []*pipelinePB.Connector{}
	for idx := range connectors {
		if _, isAirbyte := airbyteDefNames[connectors[idx].ConnectorDefinitionName]; firstProbe || !isAirbyte {
			filConnectors = append(filConnectors, connectors[idx])
		}
	}

	connectorType := "connectors"

	for _, connector := range filConnectors {

		resourcePermalink := util.ConvertUIDToResourcePermalink(connector.Uid, connectorType)

		// if user desires disconnected
		if connector.State == pipelinePB.Connector_STATE_DISCONNECTED {
			if err := s.UpdateResourceState(ctx, &controllerPB.Resource{
				ResourcePermalink: resourcePermalink,
				State: &controllerPB.Resource_ConnectorState{
					ConnectorState: pipelinePB.Connector_STATE_DISCONNECTED,
				},
			}); err != nil {
				logger.Error(err.Error())
			}
			continue
		}
		// if user desires connected
		resp, err := s.pipelinePrivateClient.CheckConnector(ctx, &pipelinePB.CheckConnectorRequest{
			Permalink: fmt.Sprintf("%s/%s", connectorType, connector.Uid),
		})

		state := pipelinePB.Connector_STATE_UNSPECIFIED
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
	}

	return nil
}
