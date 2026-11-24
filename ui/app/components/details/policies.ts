/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { policySnippetArgs } from 'core/utils/code-generators/policy';
import { terraformResourceTemplate } from 'core/utils/code-generators/terraform';
import { cliTemplate } from 'core/utils/code-generators/cli';
import type { Agent } from 'vault/agent-registry';
import type { Entity, Group } from 'vault/vault/identity';
import type { AggregatePolicy } from 'vault/utils/policy-aggregator';

interface Args {
  agent: Agent;
  entity: Entity;
  groups: Group[];
  aggregatePolicy: AggregatePolicy;
  isParentSelected: boolean;
}

export default class DetailsPoliciesComponent extends Component<Args> {
  @tracked selectedTabIndex = 1;

  @action
  onClickTab(_event: Event, index: number) {
    this.selectedTabIndex = index;
  }
  get groupPolicies() {
    return this.args.groups.flatMap((group) =>
      (group.policies ?? []).map((policy) => ({
        name: group.name,
        policy,
      }))
    );
  }

  get snippetArgs() {
    const policyName = this.args.agent?.display_name || '<policy name>';
    return policySnippetArgs(policyName, this.args.aggregatePolicy.policyString);
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
