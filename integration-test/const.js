import { uuidv4 } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';

let proto
let pHost, cHost, ctHost
let pPublicPort, pPrivatePort, cPublicPort, cPrivatePort, ctPrivatePort

export const apiGatewayMode = (__ENV.API_GATEWAY_URL && true);

if (__ENV.API_GATEWAY_PROTOCOL) {
  if (__ENV.API_GATEWAY_PROTOCOL !== "http" && __ENV.API_GATEWAY_PROTOCOL != "https") {
    fail("only allow `http` or `https` for API_GATEWAY_PROTOCOL")
  }
  proto = __ENV.API_GATEWAY_PROTOCOL
} else {
  proto = "http"
}

if (__ENV.API_GATEWAY_PROTOCOL) {
  if (__ENV.API_GATEWAY_PROTOCOL !== "http" && __ENV.API_GATEWAY_PROTOCOL != "https") {
    fail("only allow `http` or `https` for API_GATEWAY_PROTOCOL")
  }
  proto = __ENV.API_GATEWAY_PROTOCOL
} else {
  proto = "http"
}


export const pipelinePrivateHost = `http://pipeline-backend:3081`;
export const pipelinePublicHost = apiGatewayMode ? `${proto}://${__ENV.API_GATEWAY_URL}/vdp` : `http://pipeline-backend:8081`
export const pipelineGRPCPublicHost = `${__ENV.API_GATEWAY_URL}`;
export const connectorPrivateHost = `http://connector-backend:3082`;
export const connectorPublicHost = apiGatewayMode ? `${proto}://${__ENV.API_GATEWAY_URL}/vdp` : `http://connector-backend:8082`
export const connectorGRPCPrivateHost = `connector-backend:3082`;
export const connectorGRPCPublicHost = `${__ENV.API_GATEWAY_URL}`;
export const controllerPrivateHost = `http://controller-vdp:3085`;
export const controllerGRPCPrivateHost = `controller-vdp:3085`;
export const mgmtPublicHost = apiGatewayMode ? `${proto}://${__ENV.API_GATEWAY_URL}/base` : `http://mgmt-backend:8084`



export const connectorResourcePermalink = `resources/${uuidv4()}/types/connectors`

export const pipelineResourcePermalink = `resources/${uuidv4()}/types/pipelines`

export const serviceResourcePermalink = `resources/${uuidv4()}/types/services`
