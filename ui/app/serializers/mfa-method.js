import ApplicationSerializer from './application';

export default class KeymgmtKeySerializer extends ApplicationSerializer {
  normalizeItems(payload) {
    if (!payload?.data) return payload;
    if (payload.data.keys && Array.isArray(payload.data.keys)) {
      let data = payload.data.keys.map((key) => {
        let model = payload.data.key_info[key];
        model.id = key;
        return model;
      });
      return data;
    }
    Object.assign(payload, payload.data);
    delete payload.data;
    return payload;
  }
  serialize() {
    const json = super.serialize(...arguments);
    delete json.type;
    return json;
  }
}
