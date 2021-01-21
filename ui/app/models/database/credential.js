import Model, { attr } from '@ember-data/model';
// import { apiPath } from 'vault/macros/lazy-capabilities';
// import attachCapabilities from 'vault/lib/attach-capabilities';

// const ModelExport = Model.extend({
//   backend: attr('string', { readOnly: true }),
// });

// export default attachCapabilities(ModelExport, {
//   updatePath: apiPath`${'backend'}/roles/${'id'}`,
// });

export default Model.extend({
  // ARG TODO I suspect this isn't correct pattern for password
  username: attr('string'),
  password: attr('string'),
  leaseId: attr('string'),
  leaseDuration: attr('string'),
  lastVaultRotation: attr('string'),
  rotationPeriod: attr('number'),
  ttl: attr('number'),
});
