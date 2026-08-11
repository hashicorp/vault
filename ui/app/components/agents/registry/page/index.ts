/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import { service } from '@ember/service';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import {
  AGENTS_REGISTRY_CSV_FILENAME,
  CSV_COLUMNS,
  AGENT_REGISTRY_PATH,
} from 'vault/utils/constants/agent-registry';
import { generateCsv } from 'vault/utils/generate-csv';
import { dateFormat } from 'core/helpers/date-format';
import { WIZARD_ID_MAP } from 'vault/utils/constants/wizard';
import { aggregatePolicies } from 'vault/utils/policy-aggregator';

import type WizardService from 'vault/services/wizard';
import type RouterService from '@ember/routing/router-service';
import type DownloadService from 'vault/services/download';
import type ApiService from 'vault/services/api';
import type { Agent, ListAgent, RegistryCsvRow } from 'vault/agent-registry';
import type { Breadcrumb } from 'vault/app-types';
import type { AgentsRegistryRouteModel } from 'vault/routes/vault/cluster/agents/registry';
import type { AgentRegistryTableData } from '../table';
import type { Entity, Group, ListEntity } from 'vault/identity';
import type { AggregatePolicy } from 'vault/utils/policy-aggregator';

const OAUTH = 'Oauth';
const OAUTH_ACCESSOR_PATTERN = /^oauth-resource-server_(.+)_([^_]+)$/;

export type FlyoutData = {
  agent: ListAgent | Agent;
  entity?: ListEntity | Entity;
  groups: Group[];
  aggregatePolicy: AggregatePolicy;
  aliasId?: string;
};

interface Args {
  breadcrumbs: Array<Breadcrumb>;
  model: AgentsRegistryRouteModel;
}

export default class AgentsPageRegistryComponent extends Component<Args> {
  @service declare readonly download: DownloadService;
  @service declare readonly wizard: WizardService;
  @service declare readonly router: RouterService;
  @service declare readonly api: ApiService;

  wizardId = WIZARD_ID_MAP.agentRegistry;

  @tracked shouldRenderIntroModal = false;
  @tracked flyoutData: FlyoutData | null = null;
  @tracked showFlyout = false;

  @action
  openFlyout(tableData: AgentRegistryTableData) {
    this.showFlyout = true;
    this.fetchDetailData(tableData);
  }

  @action
  onFlyoutClose() {
    this.flyoutData = null;
    this.showFlyout = false;
  }

  async fetchDetailData(data: AgentRegistryTableData) {
    const { id, entity_id, entity: listEntity } = data.agent;

    const [agentResponse, entityResponse] = await Promise.allSettled([
      this.api.secrets.registrationReadById(id, AGENT_REGISTRY_PATH),
      this.api.identity.entityReadById(entity_id),
    ]);
    // default to what we have from the LIST requests in case of READ failures
    const agent = agentResponse.status === 'fulfilled' ? (agentResponse.value.data as Agent) : data.agent;
    const entity = entityResponse.status === 'fulfilled' ? (entityResponse.value.data as Entity) : listEntity;
    // fetch groups that the entity is a member of
    let groups: Group[] = [];
    if (entity && 'group_ids' in entity && entity.group_ids.length) {
      const groupRequests = entity.group_ids.map((id) => this.api.identity.groupReadById(id));
      const groupResponses = await Promise.allSettled(groupRequests);
      groups = groupResponses
        .filter((resp) => resp.status === 'fulfilled')
        .map((resp) => resp.value.data as Group);
    }
    // enrich aliases that belong to an oauth-resource-server mount with profile data
    // the mount_accessor format is: oauth-resource-server_<namespace_id>_<config_id>
    if (entity && 'group_ids' in entity && entity.aliases.length) {
      const oauthAliasRequests = entity.aliases
        .filter((alias) => OAUTH_ACCESSOR_PATTERN.test(alias.mount_accessor))
        .map((alias) => {
          // namespace_id may contain underscores; split at the last underscore so config_id is the trailing segment
          const [, namespace_id, config_id] = alias.mount_accessor.match(OAUTH_ACCESSOR_PATTERN) as [
            string,
            string,
            string,
          ];
          return this.api.sys.oauthResourceServerReadProfileById(config_id).then((profile) => {
            alias.namespace = namespace_id;
            alias.issuer_id = profile.issuer_id;
            alias.profile_name = profile.profile_name;
            alias.profile_id = profile.config_id;
            alias.mount_type = OAUTH;
          });
        });
      await Promise.allSettled(oauthAliasRequests);
    }
    // combine policy names across agent, entity and groups
    const agentPolicies =
      'ceiling_policy' in agent && agent.ceiling_policy?.length ? agent.ceiling_policy : [];
    const entityPolicies = entity?.policies?.length ? entity.policies : [];
    const groupPolicies = groups.flatMap((group) => group.policies ?? []);
    // spreading the Set will return a unique array
    const policyNames = [...new Set([...agentPolicies, ...entityPolicies, ...groupPolicies])];
    // fetch unique policies to be aggregated
    const policyRequests = policyNames.map((name) => this.api.sys.policiesReadAclPolicy(name));
    const policyResponses = await Promise.allSettled(policyRequests);
    const policies = policyResponses
      .filter((resp) => resp.status === 'fulfilled')
      .map((resp) => resp.value.policy as string);

    this.flyoutData = {
      agent,
      entity,
      groups,
      aggregatePolicy: aggregatePolicies(policies),
      aliasId: data.isAlias ? data.entityAliasId : undefined,
    };
  }

  rowsForCsv(agents: AgentsRegistryRouteModel['agents']) {
    const rows: RegistryCsvRow[] = [];
    const formatDate = (date: string) =>
      dateFormat([date, 'MMM dd, yyyy hh:mm:ss a'], {
        withTimeZone: true,
      });

    agents.forEach((agent) => {
      const { entity } = agent;

      rows.push({
        agentName: agent.display_name || '',
        agenticEntityInVault: !entity ? '/' : entity.name,
        entityOrAliasId: agent.entity_id,
        entityStatus: !entity ? '/' : entity.disabled ? 'Disabled' : 'Enabled',
        entityCreatedAt: !entity ? '/' : formatDate(entity.creation_time),
        entityUpdatedAt: !entity ? '/' : formatDate(entity.last_update_time),
      });

      (entity?.aliases || []).forEach((alias) => {
        rows.push({
          agentName: '',
          agenticEntityInVault: alias.name ? `Alias: ${alias.name}` : '/',
          entityOrAliasId: alias.id,
          entityStatus: '/',
          entityCreatedAt: formatDate(alias.creation_time),
          entityUpdatedAt: formatDate(alias.last_update_time),
        });
      });
    });

    return rows;
  }

  @action
  exportRegistrations() {
    const csvContent = generateCsv({
      rows: this.rowsForCsv(this.args.model.agents),
      columns: CSV_COLUMNS,
    });
    this.download.csv(AGENTS_REGISTRY_CSV_FILENAME, csvContent);
  }

  // wizard related
  @action
  showIntroPage() {
    // Reset the wizard dismissal state to allow re-entering the wizard
    this.wizard.reset(this.wizardId);
    this.shouldRenderIntroModal = true;
  }

  @action
  async refreshAgentsList() {
    this.router.refresh('vault.cluster.agents.registry.index');
  }

  get showContent() {
    // Show when the 1) wizard is not shown OR 2) wizard intro modal is shown
    // This ensures the wizard intro modal is shown on top of the list view and the background content is not blank behind the modal
    return !this.showWizard || (this.shouldRenderIntroModal && this.wizard.isIntroVisible(this.wizardId));
  }

  get showIntroButton() {
    return this.showContent && this.args.model.agents.length === 0;
  }

  get showWizard() {
    return !this.wizard.isDismissed(this.wizardId) && this.args.model.agents.length === 0;
  }
}
