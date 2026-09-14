/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { dateFormat } from 'core/helpers/date-format';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

import type ApiService from 'vault/services/api';
import type FlashMessagesService from 'vault/services/flash-messages';
import type RouterService from '@ember/routing/router-service';
import type { Alias, Entity } from 'vault/vault/identity';

interface Args {
  model: Entity[];
  onClickEntity?: (entity: Entity) => void;
  onClickAlias?: (entity: Entity, alias: Alias) => void;
}
const OAUTH = 'OAuth';
const OAUTH_ACCESSOR_PATTERN = /^oauth-resource-server_(.+)_([^_]+)$/;

export default class IdentityEntitiesTableComponent extends Component<Args> {
  @service declare readonly flashMessages: FlashMessagesService;
  @service declare readonly api: ApiService;
  @service declare readonly router: RouterService;

  @tracked entityToDelete = '';
  @tracked aliasToDelete = '';
  @tracked entityToDisable = '';
  @tracked filterValue = '';

  get columns() {
    return [
      { key: 'entityName', label: 'Entity name', isExpandable: true, customTableItem: true },
      { key: 'entityAliasId', label: 'Entity/Alias ID', customTableItem: true },
      { key: 'entityStatus', label: 'Entity status', customTableItem: true, width: '120px' },
      { key: 'authMethod', label: 'Auth method', customTableItem: true },
      { key: 'entityUpdatedAt', label: 'Last updated' },
      { key: 'entityCreatedAt', label: 'Created at' },
      { key: 'popupMenu', label: 'Actions', width: '80px' },
    ];
  }

  get data() {
    // return aliases as children in the data structure for nested rows and normalize keys
    const mappedData = this.args.model.map((entity) => {
      return {
        entity,
        entityName: entity.name,
        entityAliasId: entity.id,
        entityStatus: entity.disabled ? 'Disabled' : 'Enabled',
        authMethod: entity.aliases ? this.getAuthMethod(entity.aliases) : null,
        entityCreatedAt: this.formatDate(entity.creation_time),
        entityUpdatedAt: this.formatDate(entity.last_update_time),
        children:
          entity.aliases?.map((alias) => {
            return {
              isAlias: true,
              entity,
              alias,
              entityName: alias.name,
              entityAliasId: alias.id,
              authMethod: this.getAliasAuthMethod(alias),
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

  getAuthMethod(aliases: Alias[]) {
    const methods: string[] = [];
    aliases.forEach((alias) => {
      if (alias.mount_type) {
        methods.push(alias.mount_type);
      } else if (OAUTH_ACCESSOR_PATTERN.test(alias.mount_accessor)) {
        methods.push(OAUTH);
      }
    });
    if (methods.every((val) => val === methods[0])) {
      return methods[0];
    } else {
      return 'Mixed methods';
    }
  }

  getAliasAuthMethod(alias: Alias) {
    if (alias.mount_type) {
      return alias.mount_type;
    } else {
      return OAUTH_ACCESSOR_PATTERN.test(alias.mount_accessor) ? OAUTH : null;
    }
  }

  formatDate(date?: string) {
    return date ? dateFormat([date, 'MMM dd, yyyy hh:mm:ss a'], { withTimeZone: true }) : '/';
  }

  @action
  async toggleEntity(id: string, disabled: boolean) {
    try {
      await this.api.identity.entityUpdateById(id, { disabled });
      this.flashMessages.success(`Successfully ${disabled ? 'disabled' : 'enabled'} entity`);
      this.router.refresh('vault.cluster.access.identity.entities.index');
    } catch (error) {
      const { message } = await this.api.parseError(error);
      this.flashMessages.danger(`Error disabling entity: ${message}`);
    }
  }

  @action
  async deleteEntity() {
    try {
      const entity = this.entityToDelete;
      await this.api.identity.entityDeleteById(entity);
      this.flashMessages.success(`Successfully deleted entity ${entity}`);
      this.router.refresh('vault.cluster.access.identity.entities.index');
    } catch (error) {
      const { message } = await this.api.parseError(error);
      this.flashMessages.danger(`Error deleting entity ${this.entityToDelete}: ${message}`);
    }
  }

  @action
  async deleteAlias() {
    try {
      const alias = this.aliasToDelete;
      await this.api.identity.entityDeleteAliasById(alias);
      this.flashMessages.success(`Successfully deleted alias ${alias}`);
      this.router.refresh('vault.cluster.access.identity.entities.index');
    } catch (error) {
      const { message } = await this.api.parseError(error);
      this.flashMessages.danger(`Error deleting alias ${this.aliasToDelete}: ${message}`);
    }
  }
}
