import ApplicationSerializer from '../application';

export default class PkiRoleSerializer extends ApplicationSerializer {
  attrs = {
    name: { serialize: false },
  };

  serialize() {
    const json = super.serialize(...arguments);
    // empty arrays are being removed from serialized json
    // ensure that they are sent to the server, otherwise removing items will not be persisted
    json.key_usage = json.key_usage || [];
    json.ext_key_usage = json.ext_key_usage || [];
    return json;
  }
}
