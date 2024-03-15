# adcs-simulator

![Version: 0.0.8](https://img.shields.io/badge/Version-0.0.8-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 0.0.8](https://img.shields.io/badge/AppVersion-0.0.8-informational?style=flat-square)

ADCS simulator

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| simulator.arguments.dns | string | `"adcs-sim-service.adcs-issuer.svc,adcs2.example.com"` |  |
| simulator.arguments.ips | string | `"10.10.10.1,10.10.10.2"` |  |
| simulator.arguments.port | int | `8443` |  |
| simulator.clusterIssuserName | string | `"adcs-sim-adcsclusterissuer"` |  |
| simulator.configMapName | string | `"adcs-sim-configmap"` |  |
| simulator.containerPort | int | `8443` |  |
| simulator.containerSecurityContext.allowPrivilegeEscalation | bool | `false` |  |
| simulator.containerSecurityContext.capabilities.drop[0] | string | `"all"` |  |
| simulator.containerSecurityContext.readOnlyRootFilesystem | bool | `true` |  |
| simulator.deploymentName | string | `"adcs-sim-deployment"` |  |
| simulator.enabled | bool | `true` |  |
| simulator.environment.ENABLE_DEBUG | string | `"false"` |  |
| simulator.exampleCertificate.enabled | bool | `true` |  |
| simulator.exampleCertificate.name | string | `"adcs-sim-certificate"` |  |
| simulator.exampleCertificate.secretName | string | `"adcs-sim-certificate-secret"` |  |
| simulator.image.repository | string | `"djkormo/adcs-simulator"` |  |
| simulator.image.tag | string | `"0.0.8"` |  |
| simulator.livenessProbe.httpGet.path | string | `"/healthz"` |  |
| simulator.livenessProbe.httpGet.port | int | `8443` |  |
| simulator.livenessProbe.httpGet.scheme | string | `"HTTPS"` |  |
| simulator.livenessProbe.periodSeconds | int | `10` |  |
| simulator.livenessProbe.timeoutSeconds | int | `10` |  |
| simulator.podSecurityContext.runAsUser | int | `1000` |  |
| simulator.readinessProbe.httpGet.path | string | `"/readyz"` |  |
| simulator.readinessProbe.httpGet.port | int | `8443` |  |
| simulator.readinessProbe.httpGet.scheme | string | `"HTTPS"` |  |
| simulator.readinessProbe.initialDelaySeconds | int | `10` |  |
| simulator.readinessProbe.periodSeconds | int | `20` |  |
| simulator.readinessProbe.timeoutSeconds | int | `20` |  |
| simulator.resources.limits.cpu | string | `"100m"` |  |
| simulator.resources.limits.memory | string | `"500Mi"` |  |
| simulator.resources.requests.cpu | string | `"100m"` |  |
| simulator.resources.requests.memory | string | `"100Mi"` |  |
| simulator.secretCertificateName | string | `"adcs-sim-certificate-secret"` |  |
| simulator.secretName | string | `"adcs-sim-secret"` |  |
| simulator.serviceName | string | `"adcs-sim-service"` |  |
| simulator.servicePort | int | `8443` |  |

