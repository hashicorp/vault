/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';

interface Args {
  identityType: 'entity' | 'group';
  model: {
    meta: {
      total: number;
    };
  };
}

export default class EntityNavComponent extends Component<Args> {
  get description() {
    return this.args.identityType === 'entity'
      ? 'Create and manage unique identities for human and non-human identities to serve as the canonical reference ID for policies and metadata.'
      : 'Create and name logical collections of entities to simplify policy management and permission scaling across your organization.';
  }
}
