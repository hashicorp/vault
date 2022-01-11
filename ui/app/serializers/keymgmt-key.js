import ApplicationSerializer from './application';

export default class KeymgmtKeySerializer extends ApplicationSerializer {
  normalizeItems(payload) {
    let normalized = super.normalizeItems(payload);
    // Transform keys from object with number keys to array with key ids
    // Check if this is a single, list endpoint also has keys
    let lastRotated;
    if (normalized.name && normalized.keys) {
      let keys = [];
      Object.keys(normalized.keys).forEach((key, i, arr) => {
        keys.push({
          id: key,
          ...normalized.keys[key],
        });
        // Set lastRotated to the last key
        if (arr.length - 1 === i) {
          lastRotated = normalized.keys[key].creation_time;
        }
      });
      normalized.keys = keys;
    }
    return { ...normalized, last_rotated: lastRotated };
  }
}
