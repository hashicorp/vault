/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import type { Entity, Group } from 'vault/vault/identity';

interface Args {
  entity: Entity;
  groups: Group[];
}

export default class DetailsEntityComponent extends Component<Args> {
  get columns() {
    return [
      { key: 'name', label: 'In group', isExpandable: true, customTableItem: true },
      { key: 'policies', label: 'Policies', customTableItem: true },
    ];
  }
}
