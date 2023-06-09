/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import ApplicationSerializer from '../application';

export default class KvDataSerializer extends ApplicationSerializer {
  serialize(snapshot) {
    // Regardless of if CAS===true on the kv mount the UI always sends the options object with the cas version.
    const casVersion = snapshot.attr('casVersion');
    // if casVersion is 0 && the user is on the secret.edit route && does not have read permissions to data or metadata && CAS===true on the mount this will fail and we will surface the API error. We cannot intercept the failure before because we can't check if CAS is set on the mount beforehand.
    return {
      data: snapshot.attr('data'),
      options: {
        cas: casVersion,
      },
    };
  }
}
