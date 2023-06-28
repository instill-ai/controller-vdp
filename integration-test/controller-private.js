import grpc from 'k6/net/grpc';
import {
    check,
    group
} from "k6";

import * as constant from "./const.js"

const clientPrivate = new grpc.Client();
clientPrivate.load(['proto/vdp/controller/v1alpha'], 'controller_service.proto');


export function CheckConnectorResource() {
    var httpConnectorResource = {
        "resource_permalink": constant.connectorResourcePermalink,
        "connector_state": "STATE_CONNECTED"
    }

    clientPrivate.connect(constant.controllerGRPCPrivateHost, {
        plaintext: true
    });

    group("Controller API: Create connector resource state in etcd", () => {
        var resCreateConnectorHTTP = clientPrivate.invoke('vdp.controller.v1alpha.ControllerPrivateService/UpdateResource', {
            resource: httpConnectorResource
        })
        check(resCreateConnectorHTTP, {
            "vdp.controller.v1alpha.ControllerPrivateService/UpdateResource response StatusOK": (r) => r.status === grpc.StatusOK,
            "vdp.controller.v1alpha.ControllerPrivateService/UpdateResource response connectorResource resource_permalink matched": (r) => r.message.resource.resourcePermalink == httpConnectorResource.resource_permalink,
        });
    });

    group("Controller API: Get connector resource state in etcd", () => {
        var resGetConnectorHTTP = clientPrivate.invoke(`vdp.controller.v1alpha.ControllerPrivateService/GetResource`, {
            resource_permalink: httpConnectorResource.resource_permalink
        })

        check(resGetConnectorHTTP, {
            [`vdp.controller.v1alpha.ControllerPrivateService/GetResource ${httpConnectorResource.resource_permalink} response StatusOK`]: (r) => r.status === grpc.StatusOK,
            [`vdp.controller.v1alpha.ControllerPrivateService/GetResource ${httpConnectorResource.resource_permalink} response connectorResource resource_permalink matched`]: (r) => r.message.resource.resourcePermalink === httpConnectorResource.resource_permalink,
            [`vdp.controller.v1alpha.ControllerPrivateService/GetResource ${httpConnectorResource.resource_permalink} response connectorResource state matched STATE_CONNECTED`]: (r) => r.message.resource.connectorState == "STATE_CONNECTED",
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
        var resCreatePipelineHTTP = clientPrivate.invoke('vdp.controller.v1alpha.ControllerPrivateService/UpdateResource', {
            resource: httpPipelineResource
        })

        check(resCreatePipelineHTTP, {
            "vdp.controller.v1alpha.ControllerPrivateService/UpdateResource response StatusOK": (r) => r.status === grpc.StatusOK,
            "vdp.controller.v1alpha.ControllerPrivateService/UpdateResource response pipeline resource_permalink matched": (r) => r.message.resource.resourcePermalink == httpPipelineResource.resource_permalink,
        });
    });

    group("Controller API: Get pipeline resource state in etcd", () => {
        var resGetPipelineHTTP = clientPrivate.invoke(`vdp.controller.v1alpha.ControllerPrivateService/GetResource`, {
            resource_permalink: httpPipelineResource.resource_permalink
        })

        check(resGetPipelineHTTP, {
            [`vdp.controller.v1alpha.ControllerPrivateService/GetResource ${httpPipelineResource.resource_permalink} response StatusOK`]: (r) => r.status === grpc.StatusOK,
            [`vdp.controller.v1alpha.ControllerPrivateService/GetResource ${httpPipelineResource.resource_permalink} response pipeline resource_permalink matched`]: (r) => r.message.resource.resourcePermalink === httpPipelineResource.resource_permalink,
            [`vdp.controller.v1alpha.ControllerPrivateService/GetResource ${httpPipelineResource.resource_permalink} response pipeline state matched STATE_ACTIVE`]: (r) => r.message.resource.pipelineState == "STATE_ACTIVE",
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
        var resCreateServiceHTTP = clientPrivate.invoke('vdp.controller.v1alpha.ControllerPrivateService/UpdateResource', {
            resource: httpServiceResource
        })

        check(resCreateServiceHTTP, {
            "vdp.controller.v1alpha.ControllerPrivateService/UpdateResource response StatusOK": (r) => r.status === grpc.StatusOK,
            "vdp.controller.v1alpha.ControllerPrivateService/UpdateResource response service name matched": (r) => r.message.resource.name == httpServiceResource.name,
        });
    });

    group("Controller API: Get service resource state in etcd", () => {
        var resGetServiceHTTP = clientPrivate.invoke(`vdp.controller.v1alpha.ControllerPrivateService/GetResource`, {
            resource_permalink: httpServiceResource.resource_permalink
        })

        check(resGetServiceHTTP, {
            [`vdp.controller.v1alpha.ControllerPrivateService/GetResource ${httpServiceResource.resource_permalink} response StatusOK`]: (r) => r.status === grpc.StatusOK,
            [`vdp.controller.v1alpha.ControllerPrivateService/GetResource ${httpServiceResource.resource_permalink} response service name matched`]: (r) => r.message.resource.resourcePermalink === httpServiceResource.resource_permalink,
            [`vdp.controller.v1alpha.ControllerPrivateService/GetResource ${httpServiceResource.resource_permalink} response service state matched STATE_ACTIVE`]: (r) => r.message.resource.backendState == "SERVING_STATUS_SERVING",
        });
    });
}
