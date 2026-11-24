/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { AGENT_REGISTRY_PATH } from 'vault/utils/constants/agent-registry';
import { dateFormat } from 'core/helpers/date-format';

import type RouterService from '@ember/routing/router-service';
import type ApiService from 'vault/services/api';
import type FlashMessagesService from 'vault/services/flash-messages';
import type { AgentsRegistryRouteModel } from 'vault/routes/vault/cluster/agents/registry';
import type { ListAgent } from 'vault/agent-registry';

export type AgentRegistryTableData = {
  isAlias?: boolean;
  agent: ListAgent;
  agentName: string;
  entityAliasName?: string;
  entityAliasId?: string;
  entityStatus?: string;
  entityCreatedAt?: string;
  entityUpdatedAt?: string;
};

interface Args {
  model: AgentsRegistryRouteModel;
  onClick: (data: AgentRegistryTableData) => void;
}

export default class AgentRegistryTableComponent extends Component<Args> {
  @service declare readonly router: RouterService;
  @service declare readonly api: ApiService;
  @service declare readonly flashMessages: FlashMessagesService;

  @tracked agentToDelete = '';
  @tracked entityToDisable = '';
  @tracked filterValue = '';

  get columns() {
    return [
      { key: 'agentName', label: 'Agent name', isExpandable: true, customTableItem: true },
      { key: 'entityAliasName', label: 'Agentic entity in Vault', customTableItem: true },
      { key: 'entityAliasId', label: 'Entity/Alias ID', customTableItem: true },
      { key: 'entityStatus', label: 'Entity status', customTableItem: true, width: '120px' },
      { key: 'entityCreatedAt', label: 'Entity created at' },
      { key: 'entityUpdatedAt', label: 'Entity updated at' },
      { key: 'popupMenu', label: 'Actions', width: '80px' },
    ];
  }

  get data() {
    // return aliases as children in the data structure for nested rows and normalize keys
    const mappedData = this.args.model.agents.map((agent) => {
      return {
        agent,
        agentName: agent.display_name,
        entityAliasName: agent.entity?.name,
        entityAliasId: agent.entity?.id,
        entityStatus: agent.entity === undefined ? undefined : agent.entity.disabled ? 'Disabled' : 'Enabled',
        entityCreatedAt: this.formatDate(agent.entity?.creation_time),
        entityUpdatedAt: this.formatDate(agent.entity?.last_update_time),
        children:
          agent.entity?.aliases?.map((alias) => {
            return {
              isAlias: true,
              agent,
              entityAliasName: alias.name,
              entityAliasId: alias.id,
              entityCreatedAt: this.formatDate(alias.creation_time),
              entityUpdatedAt: this.formatDate(alias.last_update_time),
            };
          }) || [],
      };
    });

    if (this.filterValue) {
      // Filter data based on filterValue across all properties
      const filterLower = this.filterValue.toLowerCase();
      return mappedData.filter((item) => {
        // Check parent row properties
        const matchesParent = Object.entries(item).some(([key, value]) => {
          // Skip children array and undefined values
          if (key === 'children' || value === undefined) return false;

          // Convert value to string and check if it includes the filter
          return String(value).toLowerCase().includes(filterLower);
        });

        // Check children (alias) properties
        const matchesChildren = item.children.some((child) => {
          return Object.entries(child).some(([key, value]) => {
            // Skip isAlias flag and undefined values
            if (key === 'isAlias' || value === undefined) return false;
            // Convert value to string and check if it includes the filter
            return String(value).toLowerCase().includes(filterLower);
          });
        });

        return matchesParent || matchesChildren;
      });
    }
    // if there is nothing to filter by, return the full mapped dataset
    return mappedData;
  }

  formatDate(date?: string) {
    return date ? dateFormat([date, 'MMM dd, yyyy hh:mm:ss a'], { withTimeZone: true }) : '/';
  }

  @action
  onPageChange(page: number) {
    this.router.transitionTo('vault.cluster.agents.registry.index', {
      queryParams: { page },
    });
  }

  @action
  onPageSizeChange(pageSize: number) {
    this.router.transitionTo('vault.cluster.agents.registry.index', {
      queryParams: { page: 1, pageSize },
    });
  }

  @action
  async toggleEntity(id: string, disabled: boolean) {
    try {
      await this.api.identity.entityUpdateById(id, { disabled });
      this.flashMessages.success(`Successfully ${disabled ? 'disabled' : 'enabled'} entity`);
      this.router.refresh('vault.cluster.agents.registry.index');
    } catch (error) {
      const { message } = await this.api.parseError(error);
      this.flashMessages.danger(`Error disabling entity: ${message}`);
    }
  }

  @action
  async deleteAgent() {
    try {
      const agent = this.agentToDelete;
      await this.api.secrets.registrationDeleteByName(this.agentToDelete, AGENT_REGISTRY_PATH);
      this.flashMessages.success(`Successfully deleted agent ${agent}`);
      this.router.refresh('vault.cluster.agents.registry.index');
    } catch (error) {
      const { message } = await this.api.parseError(error);
      this.flashMessages.danger(`Error deleting agent ${this.agentToDelete}: ${message}`);
    }
  }
}
