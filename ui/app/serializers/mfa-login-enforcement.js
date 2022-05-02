import ApplicationSerializer from './application';

export default class KeymgmtProviderSerializer extends ApplicationSerializer {
  primaryKey = 'name';

  // change keys for hasMany relationships with ids in the name
  transformHasManyKeys(data, destination) {
    const keys = {
      model: ['mfa_methods', 'identity_entities', 'identity_groups'],
      server: ['mfa_method_ids', 'identity_entity_ids', 'identity_group_ids'],
    };
    keys[destination].forEach((newKey, index) => {
      const oldKey = destination === 'model' ? keys.server[index] : keys.model[index];
      delete Object.assign(data, { [newKey]: data[oldKey] })[oldKey];
    });
  }
  normalize(model, data) {
    this.transformHasManyKeys(data, 'model');
    return super.normalize(model, data);
  }
  serialize() {
    const json = super.serialize(...arguments);
    this.transformHasManyKeys(json, 'server');
    return json;
  }
}
