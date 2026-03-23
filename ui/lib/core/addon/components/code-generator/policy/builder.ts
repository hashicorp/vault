/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import Component from '@glimmer/component';
import { formatStanzas, PolicyStanza } from 'core/utils/code-generators/policy';
import { assert } from '@ember/debug';

export interface PolicyData {
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
  constructor(owner: unknown, args: Args) {
    super(owner, args);
    assert(
      '@stanzas are required and must be an array of PolicyStanza instances',
      Array.isArray(this.args.stanzas) && this.args.stanzas.every((s) => s instanceof PolicyStanza)
    );
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
