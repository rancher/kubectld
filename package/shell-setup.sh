#!/bin/bash
set -e

token=$1

mkdir -p /nonexistent
mount -t tmpfs tmpfs /nonexistent

mkdir /nonexistent/.kube
tee /nonexistent/.kube/config <<EOF
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

cp /etc/skel/.bashrc /nonexistent
echo 'PS1="> "' >> /nonexistent/.bashrc
echo . /etc/bash_completion >> /nonexistent/.bashrc
echo 'alias k="kubectl"' >> /nonexistent/.bashrc
echo 'alias ks="kubectl -n kube-system"' >> /nonexistent/.bashrc

exec su -s /bin/bash nobody
