/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationAdapter from './application';

export default ApplicationAdapter.extend({
  createRecord(store, type, snapshot) {
    const ttl = snapshot.attr('ttl');
    const roleArn = snapshot.attr('roleArn');
    const roleType = snapshot.attr('credentialType');
    let method = 'POST';
    let options;
    const data = {};
    if (roleType === 'iam_user') {
      method = 'GET';
    } else {
      if (ttl) {
        data.ttl = ttl;
      }
      if (roleType === 'assumed_role' && roleArn) {
        data.role_arn = roleArn;
      }
      options = data.ttl || data.role_arn ? { data } : {};
    }
    const role = snapshot.attr('role');
    const url = `/v1/${role.backend}/creds/${role.name}`;

    return this.ajax(url, method, options).then((response) => {
      response.id = snapshot.id;
      response.modelName = type.modelName;
      store.pushPayload(type.modelName, response);
    });
  },
});
