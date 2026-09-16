/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Form from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';

import type { Validations } from 'vault/app-types';

export interface RaftJoinFormData {
  leader_api_addr?: string;
  retry?: boolean;
  leader_ca_cert?: string;
  leader_client_cert?: string;
  leader_client_key?: string;
}

export default class RaftJoinForm extends Form<RaftJoinFormData> {
  validations: Validations = {
    leader_api_addr: [{ type: 'presence', message: "Leader API Address can't be blank." }],
  };

  formFields = [
    new FormField('leader_api_addr', 'string', { label: 'Leader API Address' }),
    new FormField('leader_ca_cert', 'string', { label: 'Leader CA Certificate', editType: 'file' }),
    new FormField('leader_client_cert', 'string', { label: 'Leader Client Certificate', editType: 'file' }),
    new FormField('leader_client_key', 'string', { label: 'Leader Client Key', editType: 'file' }),
    new FormField('retry', 'boolean', { label: 'Keep retrying to join in case of failures' }),
  ];
}
