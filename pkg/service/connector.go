package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/instill-ai/controller-vdp/internal/util"
	"github.com/instill-ai/controller-vdp/pkg/logger"

	connectorPB "github.com/instill-ai/protogen-go/vdp/connector/v1alpha"
	controllerPB "github.com/instill-ai/protogen-go/vdp/controller/v1alpha"
)

func (s *service) ProbeSourceConnectors(ctx context.Context, cancel context.CancelFunc) error {
	defer cancel()

	logger, _ := logger.GetZapLogger(ctx)

	var wg sync.WaitGroup

	resp, err := s.connectorPrivateClient.ListSourceConnectorsAdmin(ctx, &connectorPB.ListSourceConnectorsAdminRequest{})

	if err != nil {
		return err
	}

	connectors := resp.SourceConnectors
	nextPageToken := &resp.NextPageToken
	totalSize := resp.TotalSize

	for totalSize > util.DefaultPageSize {
		resp, err := s.connectorPrivateClient.ListSourceConnectorsAdmin(ctx, &connectorPB.ListSourceConnectorsAdminRequest{
			PageToken: nextPageToken,
		})

		if err != nil {
			return err
		}

		nextPageToken = &resp.NextPageToken
		totalSize -= util.DefaultPageSize
		connectors = append(connectors, resp.SourceConnectors...)
	}

	connectorType := "source-connectors"

	wg.Add(len(connectors))

	for _, connector := range connectors {

		go func(connector *connectorPB.SourceConnector) {
			defer wg.Done()

			resourcePermalink := util.ConvertUIDToResourcePermalink(connector.Uid, connectorType)

			// if user desires disconnected
			if connector.Connector.State == connectorPB.Connector_STATE_DISCONNECTED {
				if err := s.UpdateResourceState(ctx, &controllerPB.Resource{
					ResourcePermalink: resourcePermalink,
					State: &controllerPB.Resource_ConnectorState{
						ConnectorState: connectorPB.Connector_STATE_DISCONNECTED,
					},
				}); err != nil {
					logger.Error(err.Error())
					return
				}
			}
			// if user desires connected
			resp, err := s.connectorPrivateClient.CheckSourceConnector(ctx, &connectorPB.CheckSourceConnectorRequest{
				SourceConnectorPermalink: fmt.Sprintf("%s/%s", connectorType, connector.Uid),
			})
			if err != nil {
				logger.Error(err.Error())
				return
			}
			if err := s.UpdateResourceState(ctx, &controllerPB.Resource{
				ResourcePermalink: resourcePermalink,
				State: &controllerPB.Resource_ConnectorState{
					ConnectorState: resp.State,
				},
			}); err != nil {
				logger.Error(err.Error())
				return
			}

			logResp, _ := s.GetResourceState(ctx, resourcePermalink)
			logger.Info(fmt.Sprintf("[Controller] Got %v", logResp))
		}(connector)
	}

	wg.Wait()

	return nil
}

func (s *service) ProbeDestinationConnectors(ctx context.Context, cancel context.CancelFunc) error {
	defer cancel()

	logger, _ := logger.GetZapLogger(ctx)

	var wg sync.WaitGroup

	resp, err := s.connectorPrivateClient.ListDestinationConnectorsAdmin(ctx, &connectorPB.ListDestinationConnectorsAdminRequest{})

	if err != nil {
		return err
	}

	connectors := resp.DestinationConnectors
	nextPageToken := &resp.NextPageToken
	totalSize := resp.TotalSize

	for totalSize > util.DefaultPageSize {
		resp, err := s.connectorPrivateClient.ListDestinationConnectorsAdmin(ctx, &connectorPB.ListDestinationConnectorsAdminRequest{
			PageToken: nextPageToken,
		})

		if err != nil {
			return err
		}

		nextPageToken = &resp.NextPageToken
		totalSize -= util.DefaultPageSize
		connectors = append(connectors, resp.DestinationConnectors...)
	}

	connectorType := "destination-connectors"

	wg.Add(len(connectors))

	for _, connector := range connectors {

		go func(connector *connectorPB.DestinationConnector) {
			defer wg.Done()

			resourcePermalink := util.ConvertUIDToResourcePermalink(connector.Uid, connectorType)

			// if user desires disconnected
			if connector.Connector.State == connectorPB.Connector_STATE_DISCONNECTED {
				if err := s.UpdateResourceState(ctx, &controllerPB.Resource{
					ResourcePermalink: resourcePermalink,
					State: &controllerPB.Resource_ConnectorState{
						ConnectorState: connectorPB.Connector_STATE_DISCONNECTED,
					},
				}); err != nil {
					logger.Error(err.Error())
					return
				}
			}
			// if user desires connected
			resp, err := s.connectorPrivateClient.CheckDestinationConnector(ctx, &connectorPB.CheckDestinationConnectorRequest{
				DestinationConnectorPermalink: fmt.Sprintf("%s/%s", connectorType, connector.Uid),
			})
			if err != nil {
				logger.Error(err.Error())
				return
			}
			if err := s.UpdateResourceState(ctx, &controllerPB.Resource{
				ResourcePermalink: resourcePermalink,
				State: &controllerPB.Resource_ConnectorState{
					ConnectorState: resp.State,
				},
			}); err != nil {
				logger.Error(err.Error())
				return
			}
			logResp, _ := s.GetResourceState(ctx, resourcePermalink)
			logger.Info(fmt.Sprintf("[Controller] Got %v", logResp))
		}(connector)
	}

	wg.Wait()

	return nil
}
