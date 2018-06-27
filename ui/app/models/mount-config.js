import attr from 'ember-data/attr';
import Fragment from 'ember-data-model-fragments/fragment';

export default Fragment.extend({
  defaultLeaseTtl: attr({
    label: 'Default Lease TTL',
    editType: 'ttl',
  }),
  maxLeaseTtl: attr({
    label: 'Max Lease TTL',
    editType: 'ttl',
  }),
  auditNonHmacRequestKeys: attr({
    label: 'Request keys excluded from HMACing in audit',
    editType: 'stringArray',
    helpText: "Keys that will not be HMAC'd by audit devices in the request data object.",
  }),
  auditNonHmacResponseKeys: attr({
    label: 'Response keys excluded from HMACing in audit',
    editType: 'stringArray',
    helpText: "Keys that will not be HMAC'd by audit devices in the response data object.",
  }),
  listingVisibility: attr('string', {
    editType: 'boolean',
    label: 'List method when unauthenticated',
    trueValue: 'unauth',
    falseValue: '',
  }),
  passthroughRequestHeaders: attr({
    label: 'Allowed Passthrough Request headers',
    helpText: 'Headers to whitelist and pass from the request to the backend',
    editType: 'stringArray',
  }),
});
