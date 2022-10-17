import ApplicationSerializer from '../application';

export default class VersionHistorySerializer extends ApplicationSerializer {
  normalizeItems(payload) {
    if (payload.data.keys && Array.isArray(payload.data.keys)) {
      return payload.data.keys.map((key) => ({ id: key, ...payload.data.key_info[key] }));
    }
  }
}
