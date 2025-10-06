/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { HTMLElementEvent } from 'vault/forms';
import { filterIsSupported, filterTableData } from 'core/utils/client-count-utils';
import { service } from '@ember/service';
import FlagsService from 'vault/services/flags';

import type { ActivityExportData, ClientFilterTypes } from 'vault/vault/client-counts/activity-api';

// Define the base mapping to derive types from
const CLIENT_TYPE_MAP = {
  entity: 'Entity',
  'non-entity-token': 'Non-entity',
  'pki-acme': 'ACME',
  'secret-sync': 'Secret sync',
} as const;

// Dynamically derive the tab values from the mapping
type ClientListTabs = (typeof CLIENT_TYPE_MAP)[keyof typeof CLIENT_TYPE_MAP];

export interface Args {
  exportData: ActivityExportData[];
  onFilterChange: CallableFunction;
  filterQueryParams: Record<ClientFilterTypes, string>;
}

export default class ClientsClientListPageComponent extends Component<Args> {
  @service declare readonly flags: FlagsService;

  @tracked selectedTab: ClientListTabs = 'Entity';
  @tracked exportDataByTab: Record<ClientListTabs, ActivityExportData[]> = {
    Entity: [],
    'Non-entity': [],
    ACME: [],
    'Secret sync': [],
  };

  constructor(owner: unknown, args: Args) {
    super(owner, args);

    this.args.exportData.forEach((data: ActivityExportData) => {
      const tabName = CLIENT_TYPE_MAP[data.client_type];
      this.exportDataByTab[tabName].push(data);
    });

    const firstTab = Object.keys(this.exportDataByTab)[0] as ClientListTabs;
    this.selectedTab = firstTab;
  }

  get selectedTabIndex() {
    return Object.keys(this.exportDataByTab).indexOf(this.selectedTab);
  }

  // Only render tabs for whatever the export data returns
  get tabs(): ClientListTabs[] {
    return Object.keys(this.exportDataByTab) as ClientListTabs[];
  }

  @action
  handleFilter(filters: Record<ClientFilterTypes, string>) {
    this.args.onFilterChange(filters);
  }

  @action
  onClickTab(_event: HTMLElementEvent<HTMLInputElement>, idx: number) {
    const tab = this.tabs[idx];
    this.selectedTab = tab ?? this.tabs[0]!;
  }

  get filtersAreApplied() {
    return (
      Object.keys(this.args.filterQueryParams).every((f) => filterIsSupported(f)) &&
      Object.values(this.args.filterQueryParams).some((v) => !!v)
    );
  }

  // TEMPLATE HELPERS
  filterData = (dataset: ActivityExportData[]) => filterTableData(dataset, this.args.filterQueryParams);

  tableColumns(tab: ClientListTabs) {
    // all client types have values for these columns
    const defaultColumns = [
      { key: 'client_id', label: 'Client ID' },
      { key: 'client_type', label: 'Client type' },
      { key: 'namespace_path', label: 'Namespace path' },
      { key: 'namespace_id', label: 'Namespace ID' },
      {
        key: 'client_first_used_time',
        label: 'Initial usage',
        tooltip: 'When the client ID was first used in the selected billing period.',
      },
      { key: 'mount_path', label: 'Mount path' },
      { key: 'mount_type', label: 'Mount type' },
      { key: 'mount_accessor', label: 'Mount accessor' },
    ];
    // these params only have value for "entity" client types
    const entityOnly = [
      {
        key: 'entity_name',
        label: 'Entity name',
        tooltip: 'Entity name will be empty in the case of a deleted entity.',
      },
      { key: 'entity_alias_name', label: 'Entity alias name' },
      { key: 'local_entity_alias', label: 'Local entity alias' },
      { key: 'policies', label: 'Policies' },
      { key: 'entity_metadata', label: 'Entity metadata' },
      { key: 'entity_alias_metadata', label: 'Entity alias metadata' },
      { key: 'entity_alias_custom_metadata', label: 'Entity alias custom metadata' },
      { key: 'entity_group_ids', label: 'Entity group IDs' },
    ];
    return tab === 'Entity' ? [...defaultColumns, ...entityOnly] : defaultColumns;
  }
}
