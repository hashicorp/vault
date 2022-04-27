import { Factory } from 'ember-cli-mirage';

export default Factory.extend({
  auth_method_accessors: null,
  auth_method_types: null,
  identity_entity_ids: null,
  identity_group_ids: null,
  mfa_method_ids: null,
  name: null,
  namespace_id: 'root',

  afterCreate(record, server) {
    // initialize arrays and stub some data if not provided
    if (!record.name) {
      // use random string for generated name
      record.name = (Math.random() + 1).toString(36).substring(2);
    }
    if (!record.mfa_method_ids) {
      // aggregate all existing methods and choose a random one
      const methods = ['Totp', 'Duo', 'Okta', 'Pingid'].reduce((methods, type) => {
        const records = server.schema.db[`mfa${type}Methods`].where({});
        if (records.length) {
          methods.push(...records);
        }
        return methods;
      }, []);
      const method = methods.length ? methods[Math.floor(Math.random() * methods.length)] : null;
      record.mfa_method_ids = method ? [method.id] : [];
    }
    const keys = ['auth_method_accessors', 'auth_method_types', 'identity_group_ids', 'identity_entity_ids'];
    keys.forEach((key) => {
      if (!record[key]) {
        record[key] = key === 'auth_method_types' ? ['userpass'] : [];
      }
    });
  },
});
