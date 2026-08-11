/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { policySnippetArgs } from 'core/utils/code-generators/policy';
import { terraformResourceTemplate } from 'core/utils/code-generators/terraform';
import { cliTemplate } from 'core/utils/code-generators/cli';

import type { FlyoutData } from '../page';

interface Args {
  data: FlyoutData;
  isParentSelected: boolean;
}

export default class AgentRegistryDetailsPoliciesComponent extends Component<Args> {
  get groupPolicies() {
    return this.args.data.groups.flatMap((group) =>
      (group.policies ?? []).map((policy) => ({
        name: group.name,
        policy,
      }))
    );
  }

  get snippetArgs() {
    const policyName = this.args.data.agent.display_name || '<policy name>';
    return policySnippetArgs(policyName, this.args.data.aggregatePolicy.policyString);
  }

  get customTabs() {
    return [
      {
        key: 'allowed-actions',
        label: 'Allowed actions',
      },
      {
        key: 'terraform',
        label: 'Terraform Vault provider',
        snippet: terraformResourceTemplate(this.snippetArgs.terraform),
        language: 'hcl' as const,
      },
      {
        key: 'cli',
        label: 'CLI',
        snippet: cliTemplate(this.snippetArgs.cli),
        language: 'shell' as const,
      },
    ];
  }
}
