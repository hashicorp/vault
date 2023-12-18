/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import ApplicationSerializer from '../application';

export default class LdapLibrarySerializer extends ApplicationSerializer {
  primaryKey = 'name';

  attrs = {
    backend: { serialize: false },
    name: { serialize: false },
  };

  // disable_check_in_enforcement is a boolean but needs to be presented as Disabled or Enabled
  normalize(modelClass, data) {
    data.disable_check_in_enforcement = data.disable_check_in_enforcement ? 'Disabled' : 'Enabled';
    return super.normalize(modelClass, data);
  }

  serialize() {
    const json = super.serialize(...arguments);
    json.disable_check_in_enforcement = json.disable_check_in_enforcement === 'Enabled' ? false : true;
    return json;
  }
}
