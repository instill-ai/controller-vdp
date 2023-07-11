package service

import (
	"context"
	"fmt"

	"github.com/instill-ai/controller-vdp/internal/util"
	"github.com/instill-ai/controller-vdp/pkg/logger"

	connectorPB "github.com/instill-ai/protogen-go/vdp/connector/v1alpha"
	controllerPB "github.com/instill-ai/protogen-go/vdp/controller/v1alpha"
)

func (s *service) ProbeConnectors(ctx context.Context, cancel context.CancelFunc) error {
	defer cancel()

	logger, _ := logger.GetZapLogger(ctx)

	resp, err := s.connectorPrivateClient.ListConnectorsAdmin(ctx, &connectorPB.ListConnectorsAdminRequest{})

	if err != nil {
		return err
	}

	connectors := resp.Connectors
	nextPageToken := &resp.NextPageToken
	totalSize := resp.TotalSize

	for totalSize > util.DefaultPageSize {
		resp, err := s.connectorPrivateClient.ListConnectorsAdmin(ctx, &connectorPB.ListConnectorsAdminRequest{
			PageToken: nextPageToken,
		})

		if err != nil {
			return err
		}

		nextPageToken = &resp.NextPageToken
		totalSize -= util.DefaultPageSize
		connectors = append(connectors, resp.Connectors...)
	}

	connectorType := "connectors"

	for _, connector := range connectors {

		resourcePermalink := util.ConvertUIDToResourcePermalink(connector.Uid, connectorType)

		// if user desires disconnected
		if connector.State == connectorPB.Connector_STATE_DISCONNECTED {
			if err := s.UpdateResourceState(ctx, &controllerPB.Resource{
				ResourcePermalink: resourcePermalink,
				State: &controllerPB.Resource_ConnectorState{
					ConnectorState: connectorPB.Connector_STATE_DISCONNECTED,
				},
			}); err != nil {
				logger.Error(err.Error())
			}
			continue
		}
		// if user desires connected
		resp, err := s.connectorPrivateClient.CheckConnector(ctx, &connectorPB.CheckConnectorRequest{
			ConnectorPermalink: fmt.Sprintf("%s/%s", connectorType, connector.Uid),
		})

		state := connectorPB.Connector_STATE_UNSPECIFIED
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
