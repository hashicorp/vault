/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { ACL_CAPABILITIES, isAclCapability } from 'core/utils/code-generators/policy';

import type { HTMLElementEvent } from 'vault/forms';
import type { AclCapability, PolicyStanza } from 'core/utils/code-generators/policy';

interface Args {
  stanza: PolicyStanza;
  onChange: () => void;
}
export default class CodeGeneratorPolicyStanza extends Component<Args> {
  @tracked showPreview = false;

  readonly permissions = ACL_CAPABILITIES;

  hasCapability = (c: AclCapability) => this.args.stanza.capabilities.has(c);

  @action
  togglePreview() {
    this.showPreview = !this.showPreview;
  }

  @action
  setPath(event: HTMLElementEvent<HTMLInputElement>) {
    this.args.stanza.path = event.target.value;
    this.args.onChange();
  }

  @action
  setPermissions(event: HTMLElementEvent<HTMLInputElement>) {
    const { value, checked } = event.target;
    if (isAclCapability(value)) {
      const capabilities = new Set(this.args.stanza.capabilities);
      checked ? capabilities.add(value) : capabilities.delete(value);
      // Update stanza with list of capabilities
      this.args.stanza.capabilities = capabilities;
      this.args.onChange();
    }
  }
}
