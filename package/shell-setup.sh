#!/bin/bash
set -e

token=$1

mkdir -p /nonexistent
mount -t tmpfs tmpfs /nonexistent
cd /nonexistent

mkdir .kube
cat <<EOF > .kube/config
apiVersion: v1
kind: Config
clusters:
- cluster:
    api-version: v1
    certificate-authority: "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
    server: "https://$KUBERNETES_SERVICE_HOST:$KUBERNETES_SERVICE_PORT"
  name: "Default"
contexts:
- context:
    cluster: "Default"
    user: "Default"
  name: "Default"
current-context: "Default"
users:
- name: "Default"
  user:
    token: "$(cat /var/run/secrets/kubernetes.io/serviceaccount/token)"
EOF

cp /etc/skel/.bashrc .
cat >> .bashrc <<EOF
PS1="> "
. /etc/bash_completion
alias k="kubectl"
alias ks="kubectl -n kube-system"
EOF

chmod 777 .kube .bashrc
chmod 666 .kube/config

for i in $(env | cut -d "=" -f 1 | grep "CATTLE\|RANCHER"); do
    unset $i
done

exec su -s /bin/bash nobody
