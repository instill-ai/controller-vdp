import { uuidv4 } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';

let proto
let pHost, cHost, ctHost
let pPublicPort, pPrivatePort, cPublicPort, cPrivatePort, mgPrivatePort, ctPrivatePort

if (__ENV.API_GATEWAY_VDP_HOST && !__ENV.API_GATEWAY_VDP_PORT || !__ENV.API_GATEWAY_VDP_HOST && __ENV.API_GATEWAY_VDP_PORT) {
  fail("both API_GATEWAY_VDP_HOST and API_GATEWAY_VDP_PORT should be properly configured.")
}

export const apiGatewayMode = (__ENV.API_GATEWAY_VDP_HOST && __ENV.API_GATEWAY_VDP_PORT);

if (__ENV.API_GATEWAY_PROTOCOL) {
  if (__ENV.API_GATEWAY_PROTOCOL !== "http" && __ENV.API_GATEWAY_PROTOCOL != "https") {
    fail("only allow `http` or `https` for API_GATEWAY_PROTOCOL")
  }
  proto = __ENV.API_GATEWAY_PROTOCOL
} else {
  proto = "http"
}

if (apiGatewayMode) {
  // api-gateway mode
  pHost = cHost = ctHost = __ENV.API_GATEWAY_VDP_HOST
  pPrivatePort = 3081
  cPrivatePort = 3082
  ctPrivatePort = 3085
  pPublicPort = cPublicPort = 8080
} else {
  // direct microservice mode  
  pHost = "pipeline-backend"
  cHost = "connector-backend"
  ctHost = "controller-vdp"
  pPrivatePort = 3081
  cPrivatePort = 3082
  mgPrivatePort = 3084
  ctPrivatePort = 3085
  pPublicPort = 8081
  cPublicPort = 8082
}

export const connectorPublicHost = `${proto}://${cHost}:${cPublicPort}`;
export const connectorPrivateHost = `${proto}://${cHost}:${cPrivatePort}`;
export const pipelinePublicHost = `${proto}://${pHost}:${pPublicPort}`;
export const pipelinePrivateHost = `${proto}://${pHost}:${pPrivatePort}`;
export const controllerPrivateHost = `${proto}://${ctHost}:${ctPrivatePort}`;

export const controllerGRPCPrivateHost = `${ctHost}:${ctPrivatePort}`;

export const connectorResourcePermalink = `resources/${uuidv4()}/types/connectors`

export const pipelineResourcePermalink = `resources/${uuidv4()}/types/pipelines`

export const serviceResourcePermalink = `resources/${uuidv4()}/types/services`
