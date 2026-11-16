/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { baseResourceFactory } from 'vault/resources/base-factory';

interface AwsCredentialData {
  access_key?: string;
  secret_key?: string;
  security_token?: string;
  session_token?: string;
  lease_id?: string;
  renewable?: boolean;
  lease_duration?: number;
}

export interface AwsCredentialResponse {
  data?: AwsCredentialData;
  lease_id?: string;
  renewable?: boolean;
  lease_duration?: number;
}

export default class AwsCredentialResource extends baseResourceFactory<AwsCredentialData>() {
  get displayFields() {
    return [
      { name: 'access_key', label: 'Access key', value: this.access_key },
      { name: 'secret_key', label: 'Secret key', value: this.secret_key, masked: true },
      {
        name: 'security_token',
        label: 'Security token',
        value: this.session_token || this.security_token,
        masked: true,
      },
      { name: 'lease_id', label: 'Lease ID', value: this.lease_id },
      { name: 'renewable', label: 'Renewable', value: this.renewable },
      { name: 'lease_duration', label: 'Lease duration', value: this.lease_duration },
    ].filter((field) => !field.masked || field.value);
  }
  /**
   * Returns a JSON stringified object of the credential fields suitable for display/export.
   * Filters out undefined/null values and includes only the actual credentials.
   */
  get toCreds(): string {
    const props = {
      accessKey: this.access_key,
      secretKey: this.secret_key,
      securityToken: this.session_token || this.security_token,
      leaseId: this.lease_id,
    };

    const propsWithVals = Object.entries(props).reduce<Record<string, unknown>>((acc, [key, value]) => {
      if (value) {
        acc[key] = value;
      }
      return acc;
    }, {});

    return JSON.stringify(propsWithVals, null, 2);
  }
}
