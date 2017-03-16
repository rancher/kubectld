#!/bin/bash

/usr/bin/update-rancher-ssl

mkdir -p ${HOME}/.kube
cat > ${HOME}/.kube/config << EOF
apiVersion: v1
kind: Config
clusters:
- cluster:
    api-version: v1
    server: "$SERVER"
  name: "Default"
contexts:
- context:
    cluster: "Default"
  name: "Default"
current-context: "Default"
EOF
chown -R nobody:nogroup ${HOME}/.kube

helm init -c
chown -R nobody:nogroup ${HOME}/.helm

exec su -s /bin/bash nobody -p -c 'exec kubectld' 
