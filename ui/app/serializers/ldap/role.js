/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationSerializer from '../application';

export default class LdapRoleSerializer extends ApplicationSerializer {
  serialize(snapshot) {
    // remove all fields that are not relevant to specified role type
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
