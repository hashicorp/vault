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
      record.update('name', (Math.random() + 1).toString(36).substring(2));
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
      // if no methods were found create one since it is a required for login enforcements
      if (!methods.length) {
        methods.push(server.create('mfa-totp-method'));
      }
      const method = methods.length ? methods[Math.floor(Math.random() * methods.length)] : null;
      record.update('mfa_method_ids', method ? [method.id] : []);
    }
    const targets = {
      auth_method_accessors: ['auth_userpass_bb95c2b1'],
      auth_method_types: ['userpass'],
      identity_group_ids: ['34db6b52-591e-bc22-8af0-4add5e167326'],
      identity_entity_ids: ['f831667b-7392-7a1c-c0fc-33d48cb1c57d'],
    };
    for (const key in targets) {
      if (!record.key) {
        record.update(key, targets[key]);
      }
    }
  },
});
