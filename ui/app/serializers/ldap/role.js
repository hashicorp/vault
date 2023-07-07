/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import ApplicationSerializer from '../application';

export default class KubernetesConfigSerializer extends ApplicationSerializer {
  primaryKey = 'name';

  attrs = {
    backend: { serialize: false },
    name: { serialize: false },
  };

  serialize(snapshot) {
    const { fieldsForType } = snapshot.record;
    const json = super.serialize(...arguments);
    Object.keys(json).forEach((key) => {
      if (!fieldsForType.includes(key)) {
        delete json[key];
      }
    });
    return json;
  }
}
