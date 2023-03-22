/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
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
