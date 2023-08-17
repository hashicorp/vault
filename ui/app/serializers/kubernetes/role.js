/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationSerializer from '../application';

export default class KubernetesConfigSerializer extends ApplicationSerializer {
  primaryKey = 'name';

  serialize() {
    const json = super.serialize(...arguments);
    // remove backend value from payload
    delete json.backend;
    return json;
  }
}
