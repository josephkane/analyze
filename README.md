# Analyze

This repository using go 1.11 modules instead of vendoring.
In order to enable it export GO111MODULE=on  

Service configuration can be stored using JSON, TOML, YAML, HCL(v1), and Java properties config file
Also can be configured using environment variables:

*AZ_LOGGING_LEVEL `debug`  
*AZ_LOGGING_FORMATTER `TXT`  


#Analyze deployment using HELM

1. helm chart can be found on ./helm/analyze folder  
2. check that you have access to the private docker registry - https://hub.docker.com/r/supergiant/analyze/
3. generate secret of docker-registry and apply it to k8s cluster where you are planning to deploy helm chart. 
I you need deatils how to generate secret for private docker registry please check following link - https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/  
4. ```helm install --debug analyze```
5. if there is no ingress controller nginx is installed on cluster install it using command ```helm install stable/nginx-ingress```