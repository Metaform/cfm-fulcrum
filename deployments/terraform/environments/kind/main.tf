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

# CFM infrastructure
module "cfm-kind-deploy" {
  # TODO pin the module to a release tag
  source = "git::https://github.com/Metaform/connector-fabric-manager.git//deployments/terraform/environments/kind?ref=main"
}

module "cfm-testagent" {
  source = "git::https://github.com/Metaform/connector-fabric-manager.git//deployments/terraform/modules/testagent?ref=main"
  pull_policy = "Never"  # pull locally from Docker
  depends_on = [module.cfm-kind-deploy]
}

# Fulcrum infrastructure
module "postgres" {
  source = "../../modules/postgres"
}

module "fulcrum-core" {
  source = "../../modules/fulcrum-core"
  depends_on = [module.postgres]
}

# CFM Fulcrum Agent
module "cfm-agent" {
  source               = "../../modules/cfm-agent"
  enable_nodeport = true
  cfm-agent_image      = "cfm-fulcrum:latest"
  fulcrum_token        = "placeholder"
  pull_policy = "Never"  # pull locally from Docker
  fulcrum_core_service = module.fulcrum-core.fulcrum_core_service_name
  fulcrum_core_port    = module.fulcrum-core.fulcrum_core_port
  pmanager_service_url = module.cfm-kind-deploy.pmanager_internal_url
  tmanager_service_url = module.cfm-kind-deploy.tmanager_internal_url
  depends_on = [module.cfm-kind-deploy, module.cfm-testagent, module.fulcrum-core]
}

