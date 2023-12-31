import grpc from 'k6/net/grpc';
import {
  check,
  group
} from 'k6';

import * as constant from "./const.js"
import * as controller_service from './controller-private.js';
const client = new grpc.Client();
client.load(['proto/vdp/controller/v1beta'], 'controller_service.proto');

export let options = {
  setupTimeout: '10s',
  insecureSkipTLSVerify: true,
  thresholds: {
    checks: ["rate == 1.0"],
  },
};

export default function (data) {

  /*
   * Controller API - API CALLS
   */
  if (!constant.apiGatewayMode) {
    // Health check
    group("Controller API: Health check", () => {
      client.connect(constant.controllerGRPCPrivateHost, {
        plaintext: true
      });

      check(client.invoke('vdp.controller.v1beta.ControllerPrivateService/Liveness', {}), {
        'Liveness Status is OK': (r) => r && r.status === grpc.StatusOK,
        'Response status is SERVING_STATUS_SERVING': (r) => r && r.message.healthCheckResponse.status === "SERVING_STATUS_SERVING",
      });

      check(client.invoke('vdp.controller.v1beta.ControllerPrivateService/Readiness', {}), {
        'Readiness Status is OK': (r) => r && r.status === grpc.StatusOK,
        'Response status is SERVING_STATUS_SERVING': (r) => r && r.message.healthCheckResponse.status === "SERVING_STATUS_SERVING",
      });
      client.close();
    });

    controller_service.CheckConnector()
    // NOTE: we don't check pipeline state for now
    // controller_service.CheckPipelineResource()
    controller_service.CheckServiceResource()
  } else {
    console.log("No Public APIs")
  }

}

export function teardown(data) {
  if (!constant.apiGatewayMode) {
    client.connect(constant.controllerGRPCPrivateHost, {
      plaintext: true
    });
    group("Controller API: Delete all resources created by the test", () => {

      check(client.invoke(`vdp.controller.v1beta.ControllerPrivateService/DeleteResource`, {
        resource_permalink: constant.connectorPermalink
      }), {
        [`vdp.controller.v1beta.ControllerPrivateService/DeleteResource ${constant.connectorPermalink} response StatusOK`]: (r) => r.status === grpc.StatusOK,
      });

      check(client.invoke(`vdp.controller.v1beta.ControllerPrivateService/DeleteResource`, {
        resource_permalink: constant.pipelineResourcePermalink
      }), {
        [`vdp.controller.v1beta.ControllerPrivateService/DeleteResource ${constant.pipelineResourcePermalink} response StatusOK`]: (r) => r.status === grpc.StatusOK,
      });

      check(client.invoke(`vdp.controller.v1beta.ControllerPrivateService/DeleteResource`, {
        resource_permalink: constant.serviceResourcePermalink
      }), {
        [`vdp.controller.v1beta.ControllerPrivateService/DeleteResource ${constant.serviceResourcePermalink} response StatusOK`]: (r) => r.status === grpc.StatusOK,
      });
    });
    client.close();
  } else {
    console.log("No Public APIs")
  }

}
