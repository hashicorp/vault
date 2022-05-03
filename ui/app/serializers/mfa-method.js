import { assign } from '@ember/polyfills';
import ApplicationSerializer from './application';

export default class KeymgmtKeySerializer extends ApplicationSerializer {
  normalizeItems(payload) {
    if (payload.data.keys && Array.isArray(payload.data.keys)) {
      let data = payload.data.keys.map((key) => {
        let model = payload.data.key_info[key];
        model.id = key;
        return model;
      });
      return data;
    }
    assign(payload, payload.data);
    delete payload.data;
    return payload;
  }
}
