#!/bin/bash

TERRAFORM_VERSION="v0.10.0"

rm -rf vendor

govendor init

govendor fetch github.com/hashicorp/terraform/helper/schema@=$TERRAFORM_VERSION
