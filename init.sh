#!/bin/bash

TERRAFORM_VERSION="v0.8.2"

rm -rf vendor

govendor init

govendor fetch github.com/hashicorp/terraform/helper/schema@=$TERRAFORM_VERSION
