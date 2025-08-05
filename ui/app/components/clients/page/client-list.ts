/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { tracked } from '@glimmer/tracking';
import ActivityComponent from '../activity';
import { action } from '@ember/object';

export default class ClientsClientListPageComponent extends ActivityComponent {
  @tracked selectedNamespace = '';
  @tracked selectedMountPath = '';

  // TODO stubbing this action here now, but it might end up being a callback in the parent to set URL query params
  @action
  setFilter(prop: 'selectedNamespace' | 'selectedMountPath', value: string) {
    this[prop] = value;
  }

  @action
  resetFilters() {
    this.selectedNamespace = '';
    this.selectedMountPath = '';
  }

  get namespaces() {
    // TODO map over exported activity data for list of namespaces
    return ['root'];
  }

  get mountPaths() {
    // TODO map over exported activity data for list of mountPaths
    return [];
  }
}
