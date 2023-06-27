package external

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/instill-ai/controller-vdp/config"
	"github.com/instill-ai/controller-vdp/pkg/logger"

	mgmtPB "github.com/instill-ai/protogen-go/base/mgmt/v1alpha"
	connectorPB "github.com/instill-ai/protogen-go/vdp/connector/v1alpha"
	pipelinePB "github.com/instill-ai/protogen-go/vdp/pipeline/v1alpha"
	etcdv3 "go.etcd.io/etcd/client/v3"
)

// InitEtcdServiceClient initialises a TritonServiceClient instance
func InitEtcdServiceClient(ctx context.Context) *etcdv3.Client {
	logger, _ := logger.GetZapLogger(ctx)

	var clientDialOpts grpc.DialOption
	if config.Config.ConnectorBackend.HTTPS.Cert != "" && config.Config.ConnectorBackend.HTTPS.Key != "" {
		creds, err := credentials.NewServerTLSFromFile(config.Config.ConnectorBackend.HTTPS.Cert, config.Config.ConnectorBackend.HTTPS.Key)
		if err != nil {
			logger.Fatal(err.Error())
		}
		clientDialOpts = grpc.WithTransportCredentials(creds)
	} else {
		clientDialOpts = grpc.WithTransportCredentials(insecure.NewCredentials())
	}

	host := config.Config.Etcd.Host
	port := config.Config.Etcd.Port
	timeout := config.Config.Etcd.Timeout * time.Second
	// Create etcd client
	client, err := etcdv3.New(etcdv3.Config{
		Endpoints:   []string{fmt.Sprintf("%s:%s", host, port)}, //TODO: multiple nodes
		DialTimeout: timeout,
		DialOptions: []grpc.DialOption{clientDialOpts},
	})
	if err != nil {
		logger.Fatal(err.Error())
	}

	return client
}

// InitConnectorPublicServiceClient initialises a ConnectorPublicServiceClient instance
func InitConnectorPublicServiceClient(ctx context.Context) (connectorPB.ConnectorPublicServiceClient, *grpc.ClientConn) {
	logger, _ := logger.GetZapLogger(ctx)

	var clientDialOpts grpc.DialOption
	if config.Config.ConnectorBackend.HTTPS.Cert != "" && config.Config.ConnectorBackend.HTTPS.Key != "" {
		creds, err := credentials.NewServerTLSFromFile(config.Config.ConnectorBackend.HTTPS.Cert, config.Config.ConnectorBackend.HTTPS.Key)
		if err != nil {
			logger.Fatal(err.Error())
		}
		clientDialOpts = grpc.WithTransportCredentials(creds)
	} else {
		clientDialOpts = grpc.WithTransportCredentials(insecure.NewCredentials())
	}

	clientConn, err := grpc.Dial(fmt.Sprintf("%v:%v", config.Config.ConnectorBackend.Host, config.Config.ConnectorBackend.PublicPort), clientDialOpts)
	if err != nil {
		logger.Error(err.Error())
		return nil, nil
	}

	return connectorPB.NewConnectorPublicServiceClient(clientConn), clientConn
}

// InitConnectorPrivateServiceClient initialises a ConnectorPrivateServiceClient instance
func InitConnectorPrivateServiceClient(ctx context.Context) (connectorPB.ConnectorPrivateServiceClient, *grpc.ClientConn) {
	logger, _ := logger.GetZapLogger(ctx)

	var clientDialOpts grpc.DialOption
	if config.Config.ConnectorBackend.HTTPS.Cert != "" && config.Config.ConnectorBackend.HTTPS.Key != "" {
		creds, err := credentials.NewServerTLSFromFile(config.Config.ConnectorBackend.HTTPS.Cert, config.Config.ConnectorBackend.HTTPS.Key)
		if err != nil {
			logger.Fatal(err.Error())
		}
		clientDialOpts = grpc.WithTransportCredentials(creds)
	} else {
		clientDialOpts = grpc.WithTransportCredentials(insecure.NewCredentials())
	}

	clientConn, err := grpc.Dial(fmt.Sprintf("%v:%v", config.Config.ConnectorBackend.Host, config.Config.ConnectorBackend.PrivatePort), clientDialOpts)
	if err != nil {
		logger.Error(err.Error())
		return nil, nil
	}

	return connectorPB.NewConnectorPrivateServiceClient(clientConn), clientConn
}

// InitMgmtPublicServiceClient initialises a MgmtPublicServiceClient instance
func InitMgmtPublicServiceClient(ctx context.Context) (mgmtPB.MgmtPublicServiceClient, *grpc.ClientConn) {
	logger, _ := logger.GetZapLogger(ctx)

	var clientDialOpts grpc.DialOption
	var creds credentials.TransportCredentials
	var err error
	if config.Config.MgmtBackend.HTTPS.Cert != "" && config.Config.MgmtBackend.HTTPS.Key != "" {
		creds, err = credentials.NewServerTLSFromFile(config.Config.MgmtBackend.HTTPS.Cert, config.Config.MgmtBackend.HTTPS.Key)
		if err != nil {
			logger.Fatal(err.Error())
		}
		clientDialOpts = grpc.WithTransportCredentials(creds)
	} else {
		clientDialOpts = grpc.WithTransportCredentials(insecure.NewCredentials())
	}

	clientConn, err := grpc.Dial(fmt.Sprintf("%v:%v", config.Config.MgmtBackend.Host, config.Config.MgmtBackend.PublicPort), clientDialOpts)
	if err != nil {
		logger.Error(err.Error())
		return nil, nil
	}

	return mgmtPB.NewMgmtPublicServiceClient(clientConn), clientConn
}

// InitPipelinePublicServiceClient initialises a PipelinePublicServiceClient instance
func InitPipelinePublicServiceClient(ctx context.Context) (pipelinePB.PipelinePublicServiceClient, *grpc.ClientConn) {
	logger, _ := logger.GetZapLogger(ctx)

	var clientDialOpts grpc.DialOption
	var creds credentials.TransportCredentials
	var err error
	if config.Config.PipelineBackend.HTTPS.Cert != "" && config.Config.PipelineBackend.HTTPS.Key != "" {
		creds, err = credentials.NewServerTLSFromFile(config.Config.PipelineBackend.HTTPS.Cert, config.Config.PipelineBackend.HTTPS.Key)
		if err != nil {
			logger.Fatal(err.Error())
		}
		clientDialOpts = grpc.WithTransportCredentials(creds)
	} else {
		clientDialOpts = grpc.WithTransportCredentials(insecure.NewCredentials())
	}

	clientConn, err := grpc.Dial(fmt.Sprintf("%v:%v", config.Config.PipelineBackend.Host, config.Config.PipelineBackend.PublicPort), clientDialOpts)
	if err != nil {
		logger.Error(err.Error())
		return nil, nil
	}

	return pipelinePB.NewPipelinePublicServiceClient(clientConn), clientConn
}

// InitPipelinePrivateServiceClient initialises a PipelinePrivateServiceClient instance
func InitPipelinePrivateServiceClient(ctx context.Context) (pipelinePB.PipelinePrivateServiceClient, *grpc.ClientConn) {
	logger, _ := logger.GetZapLogger(ctx)

	var clientDialOpts grpc.DialOption
	var creds credentials.TransportCredentials
	var err error
	if config.Config.PipelineBackend.HTTPS.Cert != "" && config.Config.PipelineBackend.HTTPS.Key != "" {
		creds, err = credentials.NewServerTLSFromFile(config.Config.PipelineBackend.HTTPS.Cert, config.Config.PipelineBackend.HTTPS.Key)
		if err != nil {
			logger.Fatal(err.Error())
		}
		clientDialOpts = grpc.WithTransportCredentials(creds)
	} else {
		clientDialOpts = grpc.WithTransportCredentials(insecure.NewCredentials())
	}

	clientConn, err := grpc.Dial(fmt.Sprintf("%v:%v", config.Config.PipelineBackend.Host, config.Config.PipelineBackend.PrivatePort), clientDialOpts)
	if err != nil {
		logger.Error(err.Error())
		return nil, nil
	}

	return pipelinePB.NewPipelinePrivateServiceClient(clientConn), clientConn
}
