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
  appliedFilters: Record<ClientFilterTypes, string>;
}

export default class ClientsFilterToolbar extends Component<Args> {
  filterTypes = ClientFilters;

  @tracked nsLabel = '';
  @tracked mountPath = '';
  @tracked mountType = '';

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    const { nsLabel, mountPath, mountType } = this.args.appliedFilters;
    this.nsLabel = nsLabel;
    this.mountPath = mountPath;
    this.mountType = mountType;
  }

  get anyFilters() {
    return (
      Object.keys(this.args.appliedFilters).every((f) => this.supportedFilter(f)) &&
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
      this.nsLabel = '';
      this.mountPath = '';
      this.mountType = '';
    }
    // Fire callback so URL query params update when filters are cleared
    this.applyFilters();
  }

  @action
  applyFilters() {
    this.args.onFilter({
      nsLabel: this.nsLabel,
      mountPath: this.mountPath,
      mountType: this.mountType,
    });
  }

  // Helper function
  supportedFilter = (f: string): f is ClientFilterTypes =>
    Object.values(this.filterTypes).includes(f as ClientFilterTypes);
}
