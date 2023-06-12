/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import ApplicationSerializer from '../application';

export default class KvDataSerializer extends ApplicationSerializer {
  serialize(snapshot) {
    // Regardless of if CAS === true on the kv mount, the UI always sends the "options" object with the cas version.
    const casVersion = snapshot.attr('casVersion');
    return {
      data: snapshot.attr('data'),
      options: {
        cas: casVersion,
      },
    };
  }
}
