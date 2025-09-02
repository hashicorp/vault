/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

// base component for counts child routes that can be extended as needed
// contains getters that filter and extract data from activity model for use in charts

import Component from '@glimmer/component';
import { action } from '@ember/object';

import type ClientsActivityModel from 'vault/models/clients/activity';
import type { ActivityExportData, ClientFilterTypes, EntityClients } from 'core/utils/client-count-utils';

/* This component does not actually render and is the base class to house
 shared computations between the Clients::Page::Overview and Clients::Page::List components */
interface Args {
  activity: ClientsActivityModel;
  exportData: ActivityExportData[] | EntityClients[];
  onFilterChange: CallableFunction;
  filterQueryParams: Record<ClientFilterTypes, string>;
}

export default class ClientsActivityComponent extends Component<Args> {
  @action
  handleFilter(filters: Record<ClientFilterTypes, string>) {
    const { namespace_path, mount_path, mount_type } = filters;
    this.args.onFilterChange({ namespace_path, mount_path, mount_type });
  }

  @action
  resetFilters() {
    this.handleFilter({ namespace_path: '', mount_path: '', mount_type: '' });
  }
}
