/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import { service } from '@ember/service';
import Component from '@glimmer/component';
import { formatCli, writePolicy } from 'core/utils/code-generators/cli';
import { formatEot } from 'core/utils/code-generators/formatters';
import { formatStanzas, PolicyStanza } from 'core/utils/code-generators/policy';
import { terraformTemplate } from 'core/utils/code-generators/terraform';
import { assert } from '@ember/debug';

import type NamespaceService from 'vault/services/namespace';

interface PolicyData {
  policy: string;
  stanzas: PolicyStanza[];
}
interface Args {
  // Callback to pass the formatted policy and updated array of stanzas back to the parent
  onPolicyChange: (data: PolicyData) => void;
  policyName?: string;
  stanzas: PolicyStanza[];
}

export default class CodeGeneratorPolicyBuilder extends Component<Args> {
  @service declare readonly namespace: NamespaceService;

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    assert(
      '@stanzas are required and must be an array of PolicyStanza instances',
      Array.isArray(this.args.stanzas) && this.args.stanzas.every((s) => s instanceof PolicyStanza)
    );
  }

  get snippets() {
    const policyName = this.args.policyName || '<policy name>';
    const policy = formatEot(formatStanzas(this.args.stanzas));
    const options = {
      // only add namespace if we're not in root (when namespace is '')
      ...(!this.namespace.inRootNamespace ? { namespace: `"${this.namespace.path}"` } : null),
      name: `"${policyName}"`,
      policy,
    };
    return [
      {
        key: 'terraform',
        label: 'Terraform Vault Provider',
        value: terraformTemplate({ resource: 'vault_policy', options }),
        language: 'hcl',
      },
      {
        key: 'cli',
        label: 'CLI',
        value: formatCli({ command: writePolicy(policyName), content: `- ${policy}` }),
        language: 'shell',
      },
    ];
  }

  @action
  addStanza() {
    const stanzas = [...this.args.stanzas, new PolicyStanza()];
    this.updateStanzas(stanzas);
  }

  @action
  deleteStanza(stanza: PolicyStanza) {
    const remaining = [...this.args.stanzas.filter((s) => s !== stanza)];
    // Create an empty template if the only stanza was deleted
    const stanzas = remaining.length ? [...remaining] : [new PolicyStanza()];
    this.updateStanzas(stanzas);
  }

  @action
  updateStanzas(stanzas?: PolicyStanza[]) {
    const updated = stanzas ?? this.args.stanzas;
    this.args.onPolicyChange({ policy: formatStanzas(updated), stanzas: updated });
  }
}
