/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import { service } from '@ember/service';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { formatCli, writePolicy } from 'core/utils/code-generators/cli';
import { formatEot } from 'core/utils/code-generators/formatters';
import { PolicyStanza } from 'core/utils/code-generators/policy';
import { terraformTemplate } from 'core/utils/code-generators/terraform';

import type { HTMLElementEvent } from 'vault/forms';
import type NamespaceService from 'vault/services/namespace';

interface Args {
  // Callback to pass the generated policy back to the parent
  onPolicyChange: ((policy: string) => void) | undefined;
  policyName?: 'string';
}

export default class CodeGeneratorPolicyBuilder extends Component<Args> {
  @service declare readonly namespace: NamespaceService;

  snippetTypes = { terraform: 'Terraform Vault Provider', cli: 'CLI' };

  @tracked snippetType = 'terraform';
  @tracked stanzas = [new PolicyStanza()];

  get formattedPolicy() {
    return this.stanzas.map((s) => s.preview).join('\n');
  }

  get snippet() {
    const policyName = this.args.policyName || '<policy name>';
    const policy = formatEot(this.formattedPolicy);
    const options = {
      // only add namespace if we're not in root (when namespace is '')
      ...(!this.namespace.inRootNamespace ? { namespace: `"${this.namespace.path}"` } : null),
      name: `"${policyName}"`,
      policy,
    };
    switch (this.snippetType) {
      case 'terraform':
        return terraformTemplate({ resource: 'vault_policy', options });
      case 'cli':
        return formatCli({ command: writePolicy(policyName), content: `- ${policy}` });
      default:
        return '';
    }
  }

  @action
  addStanza() {
    const stanzas = [...this.stanzas, new PolicyStanza()];
    this.updateStanzas(stanzas);
  }

  @action
  deleteStanza(stanza: PolicyStanza) {
    const remaining = [...this.stanzas.filter((s) => s !== stanza)];
    // Create an empty template if the only stanza was deleted
    const stanzas = remaining.length ? [...remaining] : [new PolicyStanza()];
    this.updateStanzas(stanzas);
  }

  @action
  handleChange() {
    if (this.args.onPolicyChange) {
      this.args.onPolicyChange(this.formattedPolicy);
    }
  }

  @action
  handleRadio(event: HTMLElementEvent<HTMLInputElement>) {
    const { value } = event.target;
    this.snippetType = value;
  }

  updateStanzas(stanzas: PolicyStanza[]) {
    // Trigger an update by reassigning tracked variable
    this.stanzas = stanzas;
    this.handleChange();
  }
}
