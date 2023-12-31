package service

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/instill-ai/controller-vdp/internal/util"
	"github.com/instill-ai/controller-vdp/pkg/logger"

	controllerPB "github.com/instill-ai/protogen-go/vdp/controller/v1beta"
	pipelinePB "github.com/instill-ai/protogen-go/vdp/pipeline/v1beta"
)

func (s *service) ProbePipelines(ctx context.Context, cancel context.CancelFunc) error {
	defer cancel()

	logger, _ := logger.GetZapLogger(ctx)

	var wg sync.WaitGroup

	resp, err := s.pipelinePrivateClient.ListPipelineReleasesAdmin(ctx, &pipelinePB.ListPipelineReleasesAdminRequest{
		View: pipelinePB.Pipeline_VIEW_RECIPE.Enum(),
	})

	if err != nil {
		return err
	}

	connectors, err := s.getConnectors(ctx)
	if err != nil {
		return err
	}
	connectorNameUidMap := map[string]string{}
	for _, connector := range connectors {
		connectorNameUidMap[connector.Name] = connector.Uid
	}

	releases := resp.Releases
	nextPageToken := &resp.NextPageToken
	totalSize := resp.TotalSize

	for totalSize > util.DefaultPageSize {
		resp, err := s.pipelinePrivateClient.ListPipelineReleasesAdmin(ctx, &pipelinePB.ListPipelineReleasesAdminRequest{
			PageToken: nextPageToken,
			View:      pipelinePB.Pipeline_VIEW_RECIPE.Enum(),
		})

		if err != nil {
			return err
		}

		nextPageToken = &resp.NextPageToken
		totalSize -= util.DefaultPageSize
		releases = append(releases, resp.Releases...)
	}

	resourceType := "pipeline_releases"

	wg.Add(len(releases))

	for _, release := range releases {

		go func(release *pipelinePB.PipelineRelease) {
			defer wg.Done()

			resourcePermalink := util.ConvertUIDToResourcePermalink(release.Uid, resourceType)

			releaseResource := controllerPB.Resource{
				ResourcePermalink: resourcePermalink,
				State: &controllerPB.Resource_PipelineState{
					PipelineState: pipelinePB.State_STATE_ACTIVE,
				},
			}

			// user desires active, now check each component's state
			releaseResource.State = &controllerPB.Resource_PipelineState{PipelineState: pipelinePB.State_STATE_ERROR}

			var resources []*controllerPB.Resource

			for _, component := range release.Recipe.Components {

				if strings.HasPrefix(component.DefinitionName, "connector-definitions/") {
					if len(component.ResourceName) > 0 {
						if uid, ok := connectorNameUidMap[component.ResourceName]; ok {
							connector, err := s.GetResourceState(ctx, util.ConvertUIDToResourcePermalink(uid, "connectors"))
							if err != nil {
								resErr := s.UpdateResourceState(ctx, &releaseResource)
								if resErr != nil {
									logger.Error(fmt.Sprintf("UpdateResourceState failed for %s", component.ResourceName))
								}
								logger.Error(fmt.Sprintf("no record found for %s in etcd", component.ResourceName))
								return
							}
							resources = append(resources, connector)
						} else {
							resources = append(resources, &controllerPB.Resource{
								ResourcePermalink: "",
								State: &controllerPB.Resource_ConnectorState{
									ConnectorState: pipelinePB.Connector_STATE_ERROR,
								},
								Progress: nil,
							})
						}

					} else {
						resources = append(resources, &controllerPB.Resource{
							ResourcePermalink: "",
							State: &controllerPB.Resource_ConnectorState{
								ConnectorState: pipelinePB.Connector_STATE_ERROR,
							},
							Progress: nil,
						})
					}
				}

			}

			for _, r := range resources {
				switch v := r.State.(type) {
				case *controllerPB.Resource_ConnectorState:
					switch v.ConnectorState {
					case pipelinePB.Connector_STATE_DISCONNECTED:
						releaseResource.State = &controllerPB.Resource_PipelineState{
							PipelineState: pipelinePB.State_STATE_INACTIVE,
						}
					case pipelinePB.Connector_STATE_UNSPECIFIED:
						releaseResource.State = &controllerPB.Resource_PipelineState{
							PipelineState: pipelinePB.State_STATE_UNSPECIFIED,
						}
					case pipelinePB.Connector_STATE_ERROR:
						releaseResource.State = &controllerPB.Resource_PipelineState{
							PipelineState: pipelinePB.State_STATE_ERROR,
						}
					default:
						continue
					}
				}
				resErr := s.UpdateResourceState(ctx, &releaseResource)
				if resErr != nil {
					logger.Error(fmt.Sprintf("UpdateResourceState failed for2 %s", release.Name))
				}
				return
			}

			releaseResource.State = &controllerPB.Resource_PipelineState{
				PipelineState: pipelinePB.State_STATE_ACTIVE,
			}
			resErr := s.UpdateResourceState(ctx, &releaseResource)
			if resErr != nil {
				logger.Error(fmt.Sprintf("UpdateResourceState failed for3 %s", release.Name))
			}
		}(release)
	}

	wg.Wait()

	return nil
}
