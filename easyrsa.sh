#! /usr/bin/env bash

TEMPDIR=$(mktemp -d)

pushd "$TEMPDIR"

curl -LO https://storage.googleapis.com/kubernetes-release/easy-rsa/easy-rsa.tar.gz
tar xzf easy-rsa.tar.gz
pushd easy-rsa-master/easyrsa3

./easyrsa init-pki
./easyrsa --batch "--req-cn=hellas.default@$(date +%s)" build-ca nopass
./easyrsa --subject-alt-name="DNS:hellas.default,DNS:hellas.default.svc" --days=10000 build-server-full server nopass

popd
popd

cp "${TEMPDIR}/easy-rsa-master/easyrsa3/pki/issued/server.crt" .
cp "${TEMPDIR}/easy-rsa-master/easyrsa3/pki/private/server.key" .
