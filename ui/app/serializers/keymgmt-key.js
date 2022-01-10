import ApplicationSerializer from './application';

export default class KeymgmtKeySerializer extends ApplicationSerializer {
  normalizeItems(payload) {
    let normalized = super.normalizeItems(payload);
    // Transform keys from object with number keys to array with key ids
    // Check if this is a single, list endpoint also has keys
    if (normalized.name && normalized.keys) {
      let keys = [];
      Object.keys(normalized.keys).forEach((key) => {
        keys.push({
          id: key,
          ...normalized.keys[key],
        });
      });
      normalized.keys = keys;
    }
    return normalized;
  }
}
