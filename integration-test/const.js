import { uuidv4 } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';

let proto
let pHost, cHost, mgHost, ctHost
let pPublicPort, pPrivatePort, cPublicPort, cPrivatePort, mgPublicPort, mgPrivatePort, ctPrivatePort

if (__ENV.API_GATEWAY_HOST && !__ENV.API_GATEWAY_PORT || !__ENV.API_GATEWAY_HOST && __ENV.API_GATEWAY_PORT) {
  fail("both API_GATEWAY_HOST and API_GATEWAY_PORT should be properly configured.")
}

export const apiGatewayMode = (__ENV.API_GATEWAY_HOST && __ENV.API_GATEWAY_PORT);

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
  pHost = cHost = mgHost = ctHost = __ENV.API_GATEWAY_HOST
  pPrivatePort = 3081
  cPrivatePort = 3082
  mgPrivatePort = 3084
  ctPrivatePort = 3085
  pPublicPort = cPublicPort = mgPublicPort = 8080
} else {
  // direct microservice mode  
  mgHost = "mgmt-backend"
  pHost = "pipeline-backend"
  cHost = "connector-backend"
  ctHost = "controller-vdp"
  pPrivatePort = 3081
  cPrivatePort = 3082
  mgPrivatePort = 3084
  ctPrivatePort = 3085
  pPublicPort = 8081
  cPublicPort = 8082
  mgPublicPort = 8084
}

export const connectorPublicHost = `${proto}://${cHost}:${cPublicPort}`;
export const connectorPrivateHost = `${proto}://${cHost}:${cPrivatePort}`;
export const pipelinePublicHost = `${proto}://${pHost}:${pPublicPort}`;
export const pipelinePrivateHost = `${proto}://${pHost}:${pPrivatePort}`;
export const mgmtPublicHost = `${proto}://${mgHost}:${mgPublicPort}`;
export const mgmtPrivateHost = `${proto}://${mgHost}:${mgPrivatePort}`;
export const controllerPrivateHost = `${proto}://${ctHost}:${ctPrivatePort}`;

export const controllerGRPCPrivateHost = `${ctHost}:${ctPrivatePort}`;

export const connectorResourcePermalink = `resources/${uuidv4()}/types/connectors`

export const pipelineResourcePermalink = `resources/${uuidv4()}/types/pipelines`

export const serviceResourcePermalink = `resources/${uuidv4()}/types/services`
