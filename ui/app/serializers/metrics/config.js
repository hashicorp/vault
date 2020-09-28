import ApplicationSerializer from '../application';

export default ApplicationSerializer.extend({
  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    const normalizedPayload = {
      id: payload.id,
      queriesAvailable: payload.data.queries_available,
      defaultMonths: payload.data.default_report_months,
      retentionMonths: payload.data.retention_months,
      enabled: payload.data.enabled.includes('enabled'),
    };
    return this._super(store, primaryModelClass, normalizedPayload, id, requestType);
  },
});
