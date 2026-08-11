/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { AGENT_REGISTRY_PATH } from 'vault/utils/constants/agent-registry';
import {
  IdentityApiEntityListByIdListEnum,
  SecretsApiRegistrationListByIdListEnum,
} from '@hashicorp/vault-client-typescript';
import { ModelFrom } from 'vault/route';

import type Controller from '@ember/controller';
import type ApiService from 'vault/services/api';
import type { Breadcrumb } from 'vault/app-types';
import type { ListAgent } from 'vault/agent-registry';
import type FlashMessageService from 'vault/services/flash-messages';
import type { ListEntity } from 'vault/identity';

type RouteParams = {
  page: number;
  pageFilter: string;
  pageSize: number;
};

interface RouteController extends Controller, RouteParams {
  breadcrumbs: Array<Breadcrumb>;
}

export type AgentsRegistryRouteModel = ModelFrom<AgentsRegistryRoute>;

export default class AgentsRegistryRoute extends Route {
  @service declare readonly api: ApiService;
  @service declare readonly flashMessages: FlashMessageService;

  queryParams = {
    page: {
      refreshModel: true,
    },
    pageSize: {
      refreshModel: true,
    },
  };

  async model(params: RouteParams) {
    const { page, pageSize } = params;
    try {
      const registrationsResponse = await this.api.secrets.registrationListById(
        AGENT_REGISTRY_PATH,
        SecretsApiRegistrationListByIdListEnum.TRUE
      );
      const agents = this.api.keyInfoToArray(registrationsResponse) as ListAgent[];

      if (agents.length) {
        // fetch entities (which includes alias information) and aggregate with registration data
        try {
          const entitiesResponse = await this.api.identity.entityListById(
            IdentityApiEntityListByIdListEnum.TRUE
          );
          const entities = this.api.keyInfoToArray(entitiesResponse) as ListEntity[];

          for (const agent of agents) {
            agent.entity = entities.find((entity) => entity['id'] === agent.entity_id);
          }
        } catch {
          // ignore error which will result in missing entity/alias data
          // this will need to be handled by the table and detail components
        }
      }

      return { agents, page, pageSize };
    } catch (error: unknown) {
      const { status } = await this.api.parseError(error);

      if (status === 404) {
        return { agents: [], page, pageSize };
      }

      throw error;
    }
  }

  setupController(controller: RouteController, resolvedModel: AgentsRegistryRouteModel) {
    super.setupController(controller, resolvedModel);

    controller.breadcrumbs = [
      { label: 'Vault', route: 'vault.cluster.dashboard', icon: 'vault' },
      { label: 'Agentic security', route: 'vault.cluster.agents' },
      { label: 'Agent registry' },
    ];
  }

  resetController(controller: RouteController, isExiting: boolean) {
    if (isExiting) {
      controller.setProperties({
        page: 1,
        pageSize: 10,
      });
    }
  }
}
