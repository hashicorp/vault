/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Model from '@ember-data/model';
import { ModelValidations } from 'vault/app-types';

export default class LdapConfigModel extends Model {
  backend: string;
  binddn: string;
  bindpass: string;
  url: string;
  password_policy: string;
  starttls: boolean;
  insecure_tls: boolean;
  certificate: string;
  client_tls_cert: string;
  client_tls_key: string;
  userdn: string;
  userattr: string;
  upndomain: string;
  connection_timeout: number;
  request_timeout: number;
  validate(): ModelValidations;
}
