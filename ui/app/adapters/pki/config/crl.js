/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { encodePath } from 'vault/utils/path-encoding-helpers';
import PkiConfigAdapter from './config';

export default class PkiConfigCrlAdapter extends PkiConfigAdapter {
  namespace = 'v1';

  _url(backend) {
    return `${this.buildURL()}/${encodePath(backend)}/config/crl`;
  }
}
