/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';

import type { FlyoutData } from '../page';

interface Args {
  data: FlyoutData;
}

export default class AgentsRegistryDetailsEntityComponent extends Component<Args> {
  get columns() {
    return [
      { key: 'name', label: 'In group', isExpandable: true, customTableItem: true },
      { key: 'policies', label: 'Policies', customTableItem: true },
    ];
  }
}
