/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';

import type ApiService from 'vault/services/api';
import type VersionService from 'vault/services/version';
import type PolicyForm from 'vault/forms/policy';

interface Args {
  form: PolicyForm;
  onSave(data: PolicyForm['data']): void;
  onCancel(): void;
}

export default class PolicyTemplate extends Component<Args> {
  @service declare readonly api: ApiService;
  @service declare readonly version: VersionService;

  @tracked form = null; // form class passed to policy-form

  get policyOptions() {
    return [
      { label: 'ACL Policy', value: 'acl', isDisabled: false },
      { label: 'Role Governing Policy', value: 'rgp', isDisabled: !this.version.hasSentinel },
    ];
  }
}
