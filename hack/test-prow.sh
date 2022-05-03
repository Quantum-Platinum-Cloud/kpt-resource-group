#!/bin/bash
# Copyright 2022 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname ${BASH_SOURCE})/..
KUBEBUILDER_VERSION=2.3.1
go version

cd ${SCRIPT_ROOT}

wget https://github.com/kubernetes-sigs/kubebuilder/releases/download/v${KUBEBUILDER_VERSION}/kubebuilder_${KUBEBUILDER_VERSION}_linux_amd64.tar.gz
tar zxvf kubebuilder_${KUBEBUILDER_VERSION}_linux_amd64.tar.gz
mkdir -p /usr/local/kubebuilder/bin
cp kubebuilder_${KUBEBUILDER_VERSION}_linux_amd64/bin/* /usr/local/kubebuilder/bin

make