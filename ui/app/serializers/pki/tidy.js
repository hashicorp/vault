/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import ApplicationSerializer from '../application';

export default class PkiTidySerializer extends ApplicationSerializer {
  serialize(snapshot, tidyType) {
    const data = super.serialize(snapshot);
    if (tidyType === 'manual') {
      delete data?.enabled;
      delete data?.intervalDuration;
    }
    return data;
  }
}
