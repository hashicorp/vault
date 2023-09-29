/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { encodePath } from 'vault/utils/path-encoding-helpers';
import PkiConfigBaseAdapter from './base';

export default class PkiConfigCrlAdapter extends PkiConfigBaseAdapter {
  namespace = 'v1';

  _url(backend) {
    return `${this.buildURL()}/${encodePath(backend)}/config/crl`;
  }
}
