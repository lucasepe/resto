#!/bin/bash

KUBECONFIG=${KUBECONFIG:-~/.kube/config}
CTX=$(kubectl config current-context)
CLUSTER=$(kubectl config view -o jsonpath="{.contexts[?(@.name==\"$CTX\")].context.cluster}")
USER=$(kubectl config view -o jsonpath="{.contexts[?(@.name==\"$CTX\")].context.user}")

# Estrai valori
CLIENT_CERT=$(kubectl config view --raw -o jsonpath="{.users[?(@.name==\"$USER\")].user.client-certificate-data}")
CLIENT_KEY=$(kubectl config view --raw -o jsonpath="{.users[?(@.name==\"$USER\")].user.client-key-data}")
CA_CERT=$(kubectl config view --raw -o jsonpath="{.clusters[?(@.name==\"$CLUSTER\")].cluster.certificate-authority-data}")
SERVER_URL=$(kubectl config view -o jsonpath="{.clusters[?(@.name==\"$CLUSTER\")].cluster.server}")
TOKEN=$(kubectl config view --raw -o jsonpath="{.users[?(@.name==\"$USER\")].user.token}")
USERNAME=$(kubectl config view -o jsonpath="{.users[?(@.name==\"$USER\")].user.username}")
PASSWORD=$(kubectl config view -o jsonpath="{.users[?(@.name==\"$USER\")].user.password}")
INSECURE=$(kubectl config view -o jsonpath="{.clusters[?(@.name==\"$CLUSTER\")].cluster.insecure-skip-tls-verify}")
PROXY_URL=""  # Non sempre presente, dipende dalla tua configurazione

cat > .env <<EOF
CLIENT_CERTIFICATE_DATA=$CLIENT_CERT
CLIENT_KEY_DATA=$CLIENT_KEY
CERTIFICATE_AUTHORITY_DATA=$CA_CERT
SERVER_URL=$SERVER_URL
TOKEN=$TOKEN
USERNAME=$USERNAME
PASSWORD=$PASSWORD
INSECURE=$INSECURE
PROXY_URL=$PROXY_URL
EOF
