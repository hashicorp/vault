/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { ClientFilters, ClientFilterTypes } from 'core/utils/client-count-utils';

interface Args {
  onFilter: CallableFunction;
}

export default class ClientsFilterToolbar extends Component<Args> {
  filterTypes = ClientFilters;

  @tracked nsLabel = '';
  @tracked mountPath = '';
  @tracked mountType = '';

  get filters() {
    return Object.values(this.filterTypes);
  }

  get anyFilters() {
    return this.filters.some((f) => this[f]);
  }

  @action
  updateFilter(prop: ClientFilterTypes, value: string, close: CallableFunction) {
    this[prop] = value;
    close();
  }

  @action
  clearFilters(filterKey: ClientFilterTypes | '') {
    if (filterKey) {
      this[filterKey] = '';
    } else {
      this.nsLabel = '';
      this.mountPath = '';
      this.mountType = '';
    }
  }

  @action
  applyFilters() {
    // Send key/value pairs of filters to parent
    const filterObject = this.filters.reduce(
      (obj, filterName) => {
        const value = this[filterName];
        obj[filterName] = value;
        return obj;
      },
      {} as Record<ClientFilterTypes, string>
    );
    this.args.onFilter(filterObject);
  }
}
