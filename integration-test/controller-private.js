import grpc from 'k6/net/grpc';
import {
    check,
    group
} from "k6";

import * as constant from "./const.js"

const clientPrivate = new grpc.Client();
clientPrivate.load(['proto/vdp/controller/v1beta'], 'controller_service.proto');


export function CheckConnector() {
    var httpConnector = {
        "resource_permalink": constant.connectorPermalink,
        "connector_state": "STATE_CONNECTED"
    }

    clientPrivate.connect(constant.controllerGRPCPrivateHost, {
        plaintext: true
    });

    group("Controller API: Create connector resource state in etcd", () => {
        var resCreateConnectorHTTP = clientPrivate.invoke('vdp.controller.v1beta.ControllerPrivateService/UpdateResource', {
            resource: httpConnector
        })
        check(resCreateConnectorHTTP, {
            "vdp.controller.v1beta.ControllerPrivateService/UpdateResource response StatusOK": (r) => r.status === grpc.StatusOK,
            "vdp.controller.v1beta.ControllerPrivateService/UpdateResource response connector resource_permalink matched": (r) => r.message.resource.resourcePermalink == httpConnector.resource_permalink,
        });
    });

    group("Controller API: Get connector resource state in etcd", () => {
        var resGetConnectorHTTP = clientPrivate.invoke(`vdp.controller.v1beta.ControllerPrivateService/GetResource`, {
            resource_permalink: httpConnector.resource_permalink
        })

        check(resGetConnectorHTTP, {
            [`vdp.controller.v1beta.ControllerPrivateService/GetResource ${httpConnector.resource_permalink} response StatusOK`]: (r) => r.status === grpc.StatusOK,
            [`vdp.controller.v1beta.ControllerPrivateService/GetResource ${httpConnector.resource_permalink} response connector resource_permalink matched`]: (r) => r.message.resource.resourcePermalink === httpConnector.resource_permalink,
            [`vdp.controller.v1beta.ControllerPrivateService/GetResource ${httpConnector.resource_permalink} response connector state matched STATE_CONNECTED`]: (r) => r.message.resource.connectorState == "STATE_CONNECTED",
        });
    });
}

export function CheckPipelineResource() {
    var httpPipelineResource = {
        "resource_permalink": constant.pipelineResourcePermalink,
        "pipeline_state": "STATE_ACTIVE"
    }

    clientPrivate.connect(constant.controllerGRPCPrivateHost, {
        plaintext: true
    });

    group("Controller API: Create pipeline resource state in etcd", () => {
        var resCreatePipelineHTTP = clientPrivate.invoke('vdp.controller.v1beta.ControllerPrivateService/UpdateResource', {
            resource: httpPipelineResource
        })

        check(resCreatePipelineHTTP, {
            "vdp.controller.v1beta.ControllerPrivateService/UpdateResource response StatusOK": (r) => r.status === grpc.StatusOK,
            "vdp.controller.v1beta.ControllerPrivateService/UpdateResource response pipeline resource_permalink matched": (r) => r.message.resource.resourcePermalink == httpPipelineResource.resource_permalink,
        });
    });

    group("Controller API: Get pipeline resource state in etcd", () => {
        var resGetPipelineHTTP = clientPrivate.invoke(`vdp.controller.v1beta.ControllerPrivateService/GetResource`, {
            resource_permalink: httpPipelineResource.resource_permalink
        })

        check(resGetPipelineHTTP, {
            [`vdp.controller.v1beta.ControllerPrivateService/GetResource ${httpPipelineResource.resource_permalink} response StatusOK`]: (r) => r.status === grpc.StatusOK,
            [`vdp.controller.v1beta.ControllerPrivateService/GetResource ${httpPipelineResource.resource_permalink} response pipeline resource_permalink matched`]: (r) => r.message.resource.resourcePermalink === httpPipelineResource.resource_permalink,
            [`vdp.controller.v1beta.ControllerPrivateService/GetResource ${httpPipelineResource.resource_permalink} response pipeline state matched STATE_ACTIVE`]: (r) => r.message.resource.pipelineState == "STATE_ACTIVE",
        });
    });
}

export function CheckServiceResource() {
    var httpServiceResource = {
        "resource_permalink": constant.serviceResourcePermalink,
        "backend_state": "SERVING_STATUS_SERVING"
    }

    clientPrivate.connect(constant.controllerGRPCPrivateHost, {
        plaintext: true
    });

    group("Controller API: Create service resource state in etcd", () => {
        var resCreateServiceHTTP = clientPrivate.invoke('vdp.controller.v1beta.ControllerPrivateService/UpdateResource', {
            resource: httpServiceResource
        })

        check(resCreateServiceHTTP, {
            "vdp.controller.v1beta.ControllerPrivateService/UpdateResource response StatusOK": (r) => r.status === grpc.StatusOK,
            "vdp.controller.v1beta.ControllerPrivateService/UpdateResource response service name matched": (r) => r.message.resource.name == httpServiceResource.name,
        });
    });

    group("Controller API: Get service resource state in etcd", () => {
        var resGetServiceHTTP = clientPrivate.invoke(`vdp.controller.v1beta.ControllerPrivateService/GetResource`, {
            resource_permalink: httpServiceResource.resource_permalink
        })

        check(resGetServiceHTTP, {
            [`vdp.controller.v1beta.ControllerPrivateService/GetResource ${httpServiceResource.resource_permalink} response StatusOK`]: (r) => r.status === grpc.StatusOK,
            [`vdp.controller.v1beta.ControllerPrivateService/GetResource ${httpServiceResource.resource_permalink} response service name matched`]: (r) => r.message.resource.resourcePermalink === httpServiceResource.resource_permalink,
            [`vdp.controller.v1beta.ControllerPrivateService/GetResource ${httpServiceResource.resource_permalink} response service state matched STATE_ACTIVE`]: (r) => r.message.resource.backendState == "SERVING_STATUS_SERVING",
        });
    });
}
