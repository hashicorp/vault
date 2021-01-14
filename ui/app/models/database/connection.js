import Model, { attr } from '@ember-data/model';
import { apiPath } from 'vault/macros/lazy-capabilities';
import attachCapabilities from 'vault/lib/attach-capabilities';

const M = Model.extend({
  // ARG TODO API docs for connection https://www.vaultproject.io/api-docs/secret/databases#configure-connection

  // URL: http://127.0.0.1:8200/v1/database/config/my-db

  plugin_name: attr('string'),
});

export default attachCapabilities(M, {
  updatePath: apiPath`${'backend'}/config/${'id'}`,
});
