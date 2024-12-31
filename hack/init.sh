#! /usr/bin/env bash

set -euo pipefail

apt update
apt-get -y install curl wget lsb-release gpg vim

wget -O - https://apt.releases.hashicorp.com/gpg |  gpg --dearmor -o /usr/share/keyrings/hashicorp-archive-keyring.gpg
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(lsb_release -cs) main" |  tee /etc/apt/sources.list.d/hashicorp.list
apt update &&  apt-get -y install terraform

mkdir /terraform

cat << EOF > /terraform/main.tf
module "vpc" {
  source  = "hellas.hellas/terraform-aws-modules/terraform-aws-vpc/github"
  version = "1.63.0"

  servers = 3
}

module "bastion" {
  source  = "hellas.hellas/terraform-community-modules/tf_aws_bastion_s3_keys/github"
  version = "2.1.0"

  servers = 3
}
EOF
