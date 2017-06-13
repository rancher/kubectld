#!/bin/bash
set -e

token=$1

mkdir -p /nonexistent
mount -t tmpfs tmpfs /nonexistent
cd /nonexistent

mkdir .kube
tee .kube/config <<EOF
apiVersion: v1
kind: Config
clusters:
- cluster:
    api-version: v1
    certificate-authority: /etc/kubernetes/ssl/ca.pem
    server: "https://kubernetes.kubernetes.rancher.internal:6443"
  name: "Default"
contexts:
- context:
    cluster: "Default"
  name: "Default"
current-context: "Default"
users:
- name: "Default"
  user:
    token: "$token"
EOF

cp /etc/skel/.bashrc .
echo 'PS1="> "' >> .bashrc
echo . /etc/bash_completion >> .bashrc
echo 'alias k="kubectl"' >> .bashrc
echo 'alias ks="kubectl -n kube-system"' >> .bashrc

exec su -s /bin/bash nobody
