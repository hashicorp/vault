/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { aggregatePolicies } from 'vault/utils/policy-aggregator';

import type ApiService from 'vault/services/api';
import type { Entity, Group, ListEntity } from 'vault/vault/identity';
import type { AggregatePolicy } from 'vault/utils/policy-aggregator';
import type { Breadcrumb } from 'vault/vault/app-types';
import type { PaginatedMetadata } from 'core/utils/paginate-list';

export type EntityFlyoutData = {
  entity: Entity;
  groups: Group[];
  aggregatePolicy: AggregatePolicy;
  // set when opening the flyout directly from an alias row click
  aliasId?: string;
};

export type EntityFlyoutTab = 'policies' | 'entity' | 'alias';

export type EntitiesRouteModel = (ListEntity & { canEdit: boolean; canAddAlias: boolean })[] &
  PaginatedMetadata;

interface Args {
  model: EntitiesRouteModel;
  breadcrumbs?: Breadcrumb[];
}

export default class IdentityPageEntitiesComponent extends Component<Args> {
  @service declare readonly api: ApiService;

  @tracked entityFlyoutData: EntityFlyoutData | null = null;
  @tracked entityFlyoutInitialTab: EntityFlyoutTab = 'policies';
  // true while the entity flyout is open but data is still loading
  @tracked entityFlyoutLoading = false;

  get showEntityFlyout(): boolean {
    return this.entityFlyoutLoading || this.entityFlyoutData !== null;
  }

  @action
  async openEntityFlyout(
    listEntity: Entity,
    initialTab: EntityFlyoutTab = 'policies',
    aliasId?: string
  ): Promise<void> {
    // Show the flyout immediately with a loading state; data fills in async
    this.entityFlyoutInitialTab = initialTab;
    this.entityFlyoutLoading = true;
    this.entityFlyoutData = null;

    // Fetch the full entity to get metadata, policies, group_ids etc.
    // (the list-level entity from the table only has summary fields)
    let entity: Entity = listEntity;
    try {
      const { data } = await this.api.identity.entityReadById(listEntity.id);
      entity = data as Entity;
    } catch (_e) {
      // fall back to list-level data if the read fails
    }

    // Fetch every group the entity belongs to
    let groups: Group[] = [];
    if (entity.group_ids?.length) {
      const groupResponses = await Promise.allSettled(
        entity.group_ids.map((id) => this.api.identity.groupReadById(id))
      );
      groups = groupResponses.filter((r) => r.status === 'fulfilled').map((r) => r.value.data as Group);
    }

    // Collect all unique policy names across entity + groups
    const entityPolicies = entity.policies ?? [];
    const groupPolicies = groups.flatMap((g) => g.policies ?? []);
    const policyNames = [...new Set([...entityPolicies, ...groupPolicies])];

    // Fetch and aggregate policy strings
    const policyResponses = await Promise.allSettled(
      policyNames.map((name) => this.api.sys.policiesReadAclPolicy(name))
    );
    const policyStrings = policyResponses
      .filter((r) => r.status === 'fulfilled')
      .map((r) => r.value.policy as string);

    this.entityFlyoutLoading = false;
    this.entityFlyoutData = {
      entity,
      groups,
      aggregatePolicy: aggregatePolicies(policyStrings),
      aliasId,
    };
  }

  // Called from alias row clicks — opens the flyout at the alias tab pre-selected to the clicked alias
  @action
  openAliasFlyout(entity: Entity, alias: { id: string }): void {
    void this.openEntityFlyout(entity, 'alias', alias.id);
  }

  @action
  onEntityFlyoutClose(): void {
    this.entityFlyoutData = null;
    this.entityFlyoutLoading = false;
  }
}
