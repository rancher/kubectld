#!/bin/bash
set -e

while ! curl -s -f http://rancher-metadata/2015-12-19/stacks/Kubernetes/services/kubernetes/uuid; do
    echo Waiting for metadata
    sleep 1
done

/usr/bin/update-rancher-ssl

UUID=$(curl -s http://rancher-metadata/2015-12-19/stacks/Kubernetes/services/kubernetes/uuid)
ACTION=$(curl -s -u $CATTLE_ACCESS_KEY:$CATTLE_SECRET_KEY "$CATTLE_URL/services?uuid=$UUID" | jq -r '.data[0].actions.certificate')

if [ -n "$ACTION" ]; then
    mkdir -p /etc/kubernetes/ssl
    cd /etc/kubernetes/ssl
    curl -s -u $CATTLE_ACCESS_KEY:$CATTLE_SECRET_KEY -X POST $ACTION > certs.zip
    unzip -o certs.zip
    rm certs.zip cert.pem key.pem
fi

sleep infinity
