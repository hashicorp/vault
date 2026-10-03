/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { policySnippetArgs } from 'core/utils/code-generators/policy';

import type PolicyModel from 'vault/vault/models/policy';

interface Args {
  model: PolicyModel;
}
export default class PagePolicyShow extends Component<Args> {
  get breadcrumbs() {
    // Provide defaults so crumbs don't error as the component is torn down
    const { policyType = 'acl', id = 'policy' } = this.args.model || {};
    return [
      { label: 'Vault', route: 'vault.cluster.dashboard', icon: 'vault' },
      {
        label: `${policyType.toUpperCase()} policies`,
        route: 'vault.cluster.policies',
        model: policyType,
      },
      { label: id },
    ];
  }

  get snippetArgs() {
    return policySnippetArgs(this.args.model.id, this.args.model.policy);
  }
}
