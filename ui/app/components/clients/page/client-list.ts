/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ActivityComponent from '../activity';
import { action } from '@ember/object';
import { cached } from '@glimmer/tracking';

import type { ClientFilterTypes, MountClients } from 'core/utils/client-count-utils';

export default class ClientsClientListPageComponent extends ActivityComponent {
  @cached
  get namespaceLabels() {
    // TODO namespace list will be updated to come from the export data, not by_namespace from sys/internal/counters/activity
    return this.args.activity.byNamespace.map((n) => n.label);
  }

  @cached
  get mounts() {
    // TODO same comment here
    return this.args.activity.byNamespace.map((n) => n.mounts).flat();
  }

  @cached
  get mountPaths() {
    return [...new Set(this.mounts.map((m: MountClients) => m.label))];
  }

  @cached
  get mountTypes() {
    return [...new Set(this.mounts.map((m: MountClients) => m.mount_type))];
  }

  @action
  handleFilter(filters: Record<ClientFilterTypes, string>) {
    const { namespace_path, mount_path, mount_type } = filters;
    this.args.onFilterChange({ namespace_path, mount_path, mount_type });
  }
}
