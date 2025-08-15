/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

// base component for counts child routes that can be extended as needed
// contains getters that filter and extract data from activity model for use in charts

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { cached } from '@glimmer/tracking';
import { isSameMonth } from 'date-fns';
import { parseAPITimestamp } from 'core/utils/date-formatters';
import { service } from '@ember/service';

import type ClientsActivityModel from 'vault/models/clients/activity';
import type ClientsVersionHistoryModel from 'vault/models/clients/version-history';
import type { ClientFilterTypes } from 'core/utils/client-count-utils';
import type NamespaceService from 'vault/services/namespace';

/* This component does not actually render and is the base class to house
 shared computations between the Clients::Page::Overview and Clients::Page::List components */
interface Args {
  activity: ClientsActivityModel;
  versionHistory: ClientsVersionHistoryModel[];
  onFilterChange: CallableFunction;
  filterQueryParams: Record<ClientFilterTypes, string>;
}

export default class ClientsActivityComponent extends Component<Args> {
  @service declare readonly namespace: NamespaceService;

  @cached
  get byMonthNewClients() {
    return this.args.activity.byMonth?.map((m) => m?.new_clients) || [];
  }

  @cached
  get isCurrentMonth() {
    const { activity } = this.args;
    const current = parseAPITimestamp(activity.responseTimestamp) as Date;
    const start = parseAPITimestamp(activity.startTime) as Date;
    const end = parseAPITimestamp(activity.endTime) as Date;
    return isSameMonth(start, current) && isSameMonth(end, current);
  }

  @cached
  get isDateRange() {
    const { activity } = this.args;
    return !isSameMonth(
      parseAPITimestamp(activity.startTime) as Date,
      parseAPITimestamp(activity.endTime) as Date
    );
  }

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
