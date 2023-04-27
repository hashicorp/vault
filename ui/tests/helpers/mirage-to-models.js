/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { getContext } from '@ember/test-helpers';

export default (data) => {
  const context = getContext();
  const store = context.owner.lookup('service:store');
  const modelName = Array.isArray(data) ? data[0].modelName : data.modelName;
  const json = context.server.serializerOrRegistry.serialize(data);
  store.push(json);
  return Array.isArray(data)
    ? data.map(({ id }) => store.peekRecord(modelName, id))
    : store.peekRecord(modelName, data.id);
};
