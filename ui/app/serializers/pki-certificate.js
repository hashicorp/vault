import RESTSerializer from '@ember-data/serializer/rest';
import { isNone, isBlank } from '@ember/utils';
import { assign } from '@ember/polyfills';
import { decamelize } from '@ember/string';
// import {pki} from 'node-forge';
// const cert = pki.certificateFromPem(this.model.certificate);
// const commonName = cert.subject.getField('CN')
// const issueDate = cert.validity.notBefore
// const expiryDate = cert.validity.notAfter

export default RESTSerializer.extend({
  keyForAttribute: function(attr) {
    return decamelize(attr);
  },

  pushPayload(store, payload) {
    const transformedPayload = this.normalizeResponse(
      store,
      store.modelFor(payload.modelName),
      payload,
      payload.id,
      'findRecord'
    );
    return store.push(transformedPayload);
  },

  normalizeItems(payload) {
    if (payload.data && payload.data.keys && Array.isArray(payload.data.keys)) {
      let ret = payload.data.keys.map(key => {
        let model = {
          id_for_nav: `cert/${key}`,
          id: key,
        };
        if (payload.backend) {
          model.backend = payload.backend;
        }
        return model;
      });
      return ret;
    }
    assign(payload, payload.data);
    delete payload.data;
    return payload;
  },

  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    const responseJSON = this.normalizeItems(payload);
    const { modelName } = primaryModelClass;
    // const certMetadata = getMetadata(responseJSON.certificate), return object with
    // getMetadata is a function (make a helper...later) use forge to parse the cert
    // get info we want, and always return an object so wrap in a try
    const certMetadata = {
      common_name: 'name',
      issue_date: 'issue date',
      expiry_date: 'expiry date',
    };
    let transformedPayload = { [modelName]: { ...certMetadata, ...responseJSON } };
    console.log({ responseJSON });
    return this._super(store, primaryModelClass, transformedPayload, id, requestType);
  },

  serializeAttribute(snapshot, json, key, attributes) {
    const val = snapshot.attr(key);
    const valHasNotChanged = isNone(snapshot.changedAttributes()[key]);
    const valIsBlank = isBlank(val);
    if (attributes.options.readOnly) {
      return;
    }
    if (attributes.type === 'object' && val && Object.keys(val).length > 0 && valHasNotChanged) {
      return;
    }
    if (valIsBlank && valHasNotChanged) {
      return;
    }

    this._super(snapshot, json, key, attributes);
  },
});
