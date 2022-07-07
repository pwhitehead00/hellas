#! /usr/bin/env bash

TEMPDIR=$(mktemp -d)

NAME=$1

pushd "$TEMPDIR"

curl -s -LO https://storage.googleapis.com/kubernetes-release/easy-rsa/easy-rsa.tar.gz
tar xzf easy-rsa.tar.gz
pushd easy-rsa-master/easyrsa3

./easyrsa init-pki
./easyrsa --batch "--req-cn=${NAME}.default@$(date +%s)" build-ca nopass
./easyrsa --subject-alt-name="DNS:${NAME}.default,DNS:${NAME}.default.svc" --days=10000 build-server-full server nopass

popd
popd

echo "cert: $(base64 ${TEMPDIR}/easy-rsa-master/easyrsa3/pki/issued/server.crt)"
echo "key: $(base64 ${TEMPDIR}/easy-rsa-master/easyrsa3/pki/private/server.key)"
