/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { ClientFilters, ClientFilterTypes, filterIsSupported } from 'core/utils/client-count-utils';

interface Args {
  onFilter: CallableFunction;
  appliedFilters: Record<ClientFilterTypes, string>;
}

export default class ClientsFilterToolbar extends Component<Args> {
  filterTypes = ClientFilters;

  @tracked namespace_path = '';
  @tracked mount_path = '';
  @tracked mount_type = '';

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    const { namespace_path, mount_path, mount_type } = this.args.appliedFilters;
    this.namespace_path = namespace_path;
    this.mount_path = mount_path;
    this.mount_type = mount_type;
  }

  get anyFilters() {
    return (
      Object.keys(this.args.appliedFilters).every((f) => filterIsSupported(f)) &&
      Object.values(this.args.appliedFilters).some((v) => !!v)
    );
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
      this.namespace_path = '';
      this.mount_path = '';
      this.mount_type = '';
    }
    // Fire callback so URL query params update when filters are cleared
    this.applyFilters();
  }

  @action
  applyFilters() {
    this.args.onFilter({
      namespace_path: this.namespace_path,
      mount_path: this.mount_path,
      mount_type: this.mount_type,
    });
  }
}
