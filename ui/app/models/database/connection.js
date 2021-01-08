import Model, { attr } from '@ember-data/model';
import { computed } from '@ember/object';
import { apiPath } from 'vault/macros/lazy-capabilities';
import attachCapabilities from 'vault/lib/attach-capabilities';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';

const M = Model.extend({
  // ARG TODO API docs for connection https://www.vaultproject.io/api-docs/secret/databases#configure-connection

  // URL: http://127.0.0.1:8200/v1/database/config/my-db
  // {
  //   "data": {
  //     "allowed_roles": ["readonly"],
  //     "connection_details": {
  //       "connection_url": "{{username}}:{{password}}@tcp(127.0.0.1:3306)/",
  //       "username": "vaultuser"
  //     },
  //     "plugin_name": "mysql-database-plugin"
  //   }
  // }
  plugin_name: attr('string'),
});

export default attachCapabilities(M, {
  updatePath: apiPath`${'backend'}/config/${'id'}`,
});
