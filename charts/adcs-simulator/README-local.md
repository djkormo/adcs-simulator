

helm template adcs-simulator charts/adcs-simulator -n adcs-issuer --values charts/adcs-simulator/values.yaml

helm template adcs-issuer charts/adcs-simulator -n adcs-issuer --values charts/adcs-simulator/values.yaml > adcs-simulator-all.yaml

kubectl -n adcs-issuer apply -f adcs-simulator-all.yaml --dry-run=server

kubectl -n adcs-issuer get pod

kubectl -n adcs-issuer logs deploy/adcs-simulator-controller-manager

kubectl -n adcs-issuer logs deploy/adcs-simulator-deployment

kubectl -n adcs-issuer delete -f adcs-simulator-all.yaml


helm-docs --sort-values-order file 

