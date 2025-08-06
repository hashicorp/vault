/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { tracked } from '@glimmer/tracking';
import ActivityComponent from '../activity';
import { action } from '@ember/object';
import { ACTIVITY_EXPORT_COLUMNS } from 'core/utils/client-count-utils';
import { toLabel } from 'core/helpers/to-label';

export default class ClientsClientListPageComponent extends ActivityComponent {
  @tracked selectedNamespace = '';
  @tracked selectedMountPath = '';

  records: Record<string, unknown>[] = [];

  constructor(owner: unknown, args: any) {
    super(owner, args);
    const non_entity = (i: number) => ({
      id: i,
      local_entity_alias: false,
      client_id: `${i}g0A6QFkmfBoCzvwgAWlrwXxqmnxezsz9P+RmbKWtbTM=`,
      client_type: 'non-entity-token',
      namespace_id: 'root',
      namespace_path: '',
      mount_accessor: 'auth_token_8a563973',
      mount_type: 'token',
      mount_path: 'auth/token/',
      timestamp: '2025-07-29T00:48:55Z',
    });
    const entity = (i: number) => ({
      id: i,
      entity_name: 'entity_106e123d',
      entity_alias_name: 'approver',
      local_entity_alias: false,
      client_id: `${i}-f4fe6aac-2a6f-5c0f-1047-174caa66fdac`,
      client_type: 'entity',
      namespace_id: 'root',
      namespace_path: '',
      mount_accessor: 'auth_userpass_e759024e',
      mount_type: 'userpass',
      mount_path: 'auth/auto/eng/core/auth/core-gh-auth/',
      timestamp: '2025-07-29T20:01:50Z',
    });
    this.records = Array.from({ length: 2500 }, (_, i) => {
      return i % 2 === 0 ? entity(i) : non_entity(i);
    });
    // eslint-disable-next-line
    console.log(this.records);
  }

  get columns() {
    const notSortable = ['client_id', 'local_entity_alias'];
    const isSortable = (key: string) => !notSortable.includes(key);
    const upperCaseID = (str: string) => str.replace(/\bid\b/gi, 'ID');

    return ['id', ...ACTIVITY_EXPORT_COLUMNS.shared].map((key) => ({
      key,
      label: upperCaseID(toLabel([key])),
      isSortable: isSortable(key),
    }));
  }

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
