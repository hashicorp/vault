/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Form from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';

import type { SystemWriteConfigControlGroupRequest } from '@hashicorp/vault-client-typescript';

export default class ControlGroupConfigForm extends Form<SystemWriteConfigControlGroupRequest> {
  formFields = [
    new FormField('max_ttl', undefined, {
      label: 'Maximum TTL',
      editType: 'ttl',
    }),
  ];
}
