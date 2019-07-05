import BaseAdapter from './base';
import { decamelize } from '@ember/string';
import { getProperties } from '@ember/object';

export default BaseAdapter.extend({
  createRecord(store, type, snapshot) {
    let name = snapshot.id || snapshot.attr('name');
    let url = this._url(
      type.modelName,
      {
        backend: snapshot.record.backend,
        scope: snapshot.record.scope,
      },
      name
    );
    return this.ajax(url, 'POST', { data: this.serialize(snapshot) }).then(() => {
      return {
        id: name,
        name,
        backend: snapshot.record.backend,
        scope: snapshot.record.scope,
      };
    });
  },

  deleteRecord(store, type, snapshot) {
    let name = snapshot.id || snapshot.attr('name');
    let url = this._url(
      type.modelName,
      {
        backend: snapshot.record.backend,
        scope: snapshot.record.scope,
      },
      name
    );
    return this.ajax(url, 'DELETE');
  },

  serialize(snapshot) {
    // the endpoint here won't allow sending `operation_all` and `operation_none` at the same time or with
    // other values, so we manually check for them and send an abbreviated object
    let json = snapshot.serialize();
    let keys = snapshot.record.nonOperationFields.map(decamelize);
    let nonOperationFields = getProperties(json, keys);
    for (let field in nonOperationFields) {
      if (nonOperationFields[field] == null) {
        delete nonOperationFields[field];
      }
    }
    if (json.operation_all) {
      return {
        operation_all: true,
        ...nonOperationFields,
      };
    }
    if (json.operation_none) {
      return {
        operation_none: true,
        ...nonOperationFields,
      };
    }
    delete json.operation_none;
    delete json.operation_all;
    return json;
  },

  updateRecord() {
    return this.createRecord(...arguments);
  },
});
