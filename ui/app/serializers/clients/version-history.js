import ApplicationSerializer from '../application';

export default class VersionHistorySerializer extends ApplicationSerializer {
  primaryKey = 'version';

  normalizeItems(payload) {
    if (payload.data.keys && Array.isArray(payload.data.keys)) {
      return payload.data.keys.map((key) => ({ version: key, ...payload.data.key_info[key] }));
    }
  }
}
