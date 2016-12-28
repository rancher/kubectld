#!/bin/bash

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

/usr/bin/update-rancher-ssl
helm init -c

exec kubectld --server=$SERVER --listen=$LISTEN
