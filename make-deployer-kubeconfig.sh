#!/usr/bin/env bash
# Builds a kubeconfig for the github-deployer ServiceAccount.
# Run on your laptop, where your admin kubectl already works.
set -euo pipefail

SERVER=https://140.238.231.41:6443
NS=go-lab
SECRET=github-deployer-token

TOKEN=$(kubectl -n "$NS" get secret "$SECRET" -o jsonpath='{.data.token}' | base64 --decode)
CA=$(kubectl -n "$NS" get secret "$SECRET" -o jsonpath='{.data.ca\.crt}')

cat > deployer-kubeconfig.yaml <<CONFIG
apiVersion: v1
kind: Config
clusters:
  - name: k3s
    cluster:
      server: $SERVER
      certificate-authority-data: $CA
users:
  - name: github-deployer
    user:
      token: $TOKEN
contexts:
  - name: deployer
    context:
      cluster: k3s
      user: github-deployer
      namespace: $NS
current-context: deployer
CONFIG

chmod 600 deployer-kubeconfig.yaml
echo "Wrote deployer-kubeconfig.yaml"
