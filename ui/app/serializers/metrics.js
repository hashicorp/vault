import ApplicationSerializer from './application';

export default ApplicationSerializer.extend({
  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    let totalTokens = payload.data.counters.service_tokens ? payload.data.counters.service_tokens.total : 0;
    let totalEntities = payload.data.counters.entities ? payload.data.counters.entities.total : 0;

    let normalizedPayload = {
      id: payload.id,
      total_tokens: totalTokens,
      total_entities: totalEntities,
    };

    return this._super(store, primaryModelClass, normalizedPayload, id, requestType);
  },
});
