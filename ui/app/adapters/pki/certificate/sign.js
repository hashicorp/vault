/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { encodePath } from 'vault/utils/path-encoding-helpers';
import PkiCertificateBaseAdapter from './base';

export default class PkiCertificateSignAdapter extends PkiCertificateBaseAdapter {
  urlForCreateRecord(modelName, snapshot) {
    const { role, backend } = snapshot.record;
    if (!role || !backend) {
      throw new Error('URL for create record is missing required attributes');
    }
    return `${this.buildURL()}/${encodePath(backend)}/sign/${encodePath(role)}`;
  }
}
