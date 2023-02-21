import ApplicationSerializer from '../application';

export default ApplicationSerializer.extend({
  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    if (!payload.data) {
      // ARG TODO RETURN
      // CBS TODO: Remove this if block once API is published
      return this._super(store, primaryModelClass, payload, id, requestType);
    }
    const normalizedPayload = {
      id: payload.id,
      data: {
        ...payload.data,
        enabled: payload.data.enabled?.includes('enable') ? 'On' : 'Off',
      },
    };
    return this._super(store, primaryModelClass, normalizedPayload, id, requestType);
  },

  serialize() {
    const json = this._super(...arguments);
    if (json.enabled === 'On' || json.enabled === 'Off') {
      const oldEnabled = json.enabled;
      json.enabled = oldEnabled === 'On' ? 'enable' : 'disable';
    }
    json.retention_months = parseInt(json.retention_months, 10);
    if (isNaN(json.retention_months)) {
      throw new Error('Invalid number value');
    }
    delete json.queries_available;
    return json;
  },
});
