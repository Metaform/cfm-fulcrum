#  Copyright (c) 2025 Metaform Systems, Inc
#
#  This program and the accompanying materials are made available under the
#  terms of the Apache License, Version 2.0 which is available at
#  https://www.apache.org/licenses/LICENSE-2.0
#
#  SPDX-License-Identifier: Apache-2.0
#
#  Contributors:
#       Metaform Systems, Inc. - initial API and implementation
#

provider "kubernetes" {
  config_path    = "~/.kube/config"
  config_context = var.kubeconfig_path
}

module "cfm-kind-deploy" {
  # TODO pin the module to a release tag
  source = "git::https://github.com/Metaform/connector-fabric-manager.git//deployments/terraform/environments/kind?ref=main"
}

module "cfm-testagent" {
  source = "git::https://github.com/Metaform/connector-fabric-manager.git//deployments/terraform/modules/testagent?ref=main"
  pull_policy     = "Never"  # pull locally from Docker
  depends_on = [module.cfm-kind-deploy]
}

module "cfm-agent" {
  source          = "../../modules/cfm-agent"
  cfm-agent_image = "cfm-fulcrum:latest"
  fulcrum_token   = "123"
  pull_policy     = "Never"  # pull locally from Docker
  enable_nodeport = true

  depends_on = [module.cfm-kind-deploy, module.cfm-testagent]
}

