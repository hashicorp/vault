/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';

interface Args {
  onFilter: CallableFunction;
}

enum SupportedFilterTypes {
  NAMESPACE = 'namespace_path',
  MOUNT_PATH = 'mount_path',
  MOUNT_TYPE = 'mount_type',
}

export default class ClientsFilterToolbar extends Component<Args> {
  filterTypes = SupportedFilterTypes;

  @tracked namespace_path = '';
  @tracked mount_path = '';
  @tracked mount_type = '';

  get filters() {
    return Object.values(this.filterTypes);
  }

  get anyFilters() {
    return this.filters.some((f) => this[f]);
  }

  @action
  updateFilter(prop: SupportedFilterTypes, value: string, close: CallableFunction) {
    this[prop] = value;
    close();
  }

  @action
  clearFilters(filterKey: SupportedFilterTypes | '') {
    if (filterKey) {
      this[filterKey] = '';
    } else {
      this.namespace_path = '';
      this.mount_path = '';
      this.mount_type = '';
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
      {} as Record<SupportedFilterTypes, string>
    );
    this.args.onFilter(filterObject);
  }
}
