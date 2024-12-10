/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import RESTSerializer from '@ember-data/serializer/rest';
import { AVAILABLE_PLUGIN_TYPES } from '../../utils/model-helpers/database-helpers';

export default RESTSerializer.extend({
  primaryKey: 'name',

  serializeAttribute(snapshot, json, key, attributes) {
    // Don't send values that are undefined
    if (undefined !== snapshot.attr(key)) {
      this._super(snapshot, json, key, attributes);
    }
  },

  normalizeSecrets(payload) {
    if (payload.data.keys && Array.isArray(payload.data.keys)) {
      const connections = payload.data.keys.map((secret) => ({ name: secret, backend: payload.backend }));
      return connections;
    }
    // Query single record response:
    const response = {
      id: payload.id,
      name: payload.id,
      backend: payload.backend,
      ...payload.data,
      ...payload.data.connection_details,
    };

    // connection_details are spread above into the main body of response so we can remove redundant data
    delete response.connection_details;
    if (response?.connection_url) {
      // this url can include interpolated data, such as: "{{username}}/{{password}}@localhost:1521/OraDoc.localhost"
      // these curly brackets are returned by the API encoded: "%7B%7Busername%7D%7D/%7B%7Bpassword%7D%7D@localhost:1521/OraDoc.localhost"
      // we decode here so the UI displays and submits the url in the correct format
      response.connection_url = decodeURI(response.connection_url);
    }

    if (payload.data.root_credentials_rotate_statements) {
      response.root_rotation_statements = payload.data.root_credentials_rotate_statements;
    }
    return response;
  },

  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    const nullResponses = ['updateRecord', 'createRecord', 'deleteRecord'];
    const connections = nullResponses.includes(requestType)
      ? { name: payload.data.name, backend: payload.data.backend }
      : this.normalizeSecrets(payload);
    const { modelName } = primaryModelClass;
    let transformedPayload = { [modelName]: connections };
    if (requestType === 'queryRecord') {
      // comes back as object anyway
      transformedPayload = { [modelName]: { id, ...connections } };
    }
    return this._super(store, primaryModelClass, transformedPayload, id, requestType);
  },

  serialize(snapshot, requestType) {
    const data = this._super(snapshot, requestType);
    if (!data.plugin_name) {
      return data;
    }
    const pluginType = AVAILABLE_PLUGIN_TYPES.find((plugin) => plugin.value === data.plugin_name);
    if (!pluginType) {
      return data;
    }
    const pluginAttributes = pluginType.fields.map((fields) => fields.attr).concat('backend');

    // filter data to only allow plugin specific attrs
    const allowedAttributes = Object.keys(data).filter((dataAttrs) => pluginAttributes.includes(dataAttrs));
    for (const key in data) {
      // All connections allow allowed_roles but it's not shown on the form
      if (key !== 'allowed_roles' && !allowedAttributes.includes(key)) {
        delete data[key];
      }
    }
    return data;
  },
});
