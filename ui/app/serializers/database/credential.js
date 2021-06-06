import RESTSerializer from '@ember-data/serializer/rest';

export default RESTSerializer.extend({
  primaryKey: 'request_id',

  normalizePayload(payload) {
    if (payload.data) {
      const credentials = {
        request_id: payload.request_id,
        username: payload.data.username,
        password: payload.data.password,
        leaseId: payload.lease_id,
        leaseDuration: payload.lease_duration,
        lastVaultRotation: payload.data.last_vault_rotation,
        rotationPeriod: payload.data.rotation_period,
        ttl: payload.data.ttl,
      };
      return credentials;
    }
  },

  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    const credentials = this.normalizePayload(payload);
    const { modelName } = primaryModelClass;
    let transformedPayload = { [modelName]: credentials };

    return this._super(store, primaryModelClass, transformedPayload, id, requestType);
  },
});
