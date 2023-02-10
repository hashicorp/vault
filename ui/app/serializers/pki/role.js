import ApplicationSerializer from '../application';

export default class PkiRoleSerializer extends ApplicationSerializer {
  attrs = {
    name: { serialize: false },
  };

  serialize() {
    const json = super.serialize(...arguments);
    // attributes with empty arrays are stripped from serialized json
    // but an empty list is acceptable for key_usage to specify no default constraints
    // intercepting here to ensure an empty array persists (the backend assumes default values)
    json.key_usage = json.key_usage || [];
    return json;
  }
}
