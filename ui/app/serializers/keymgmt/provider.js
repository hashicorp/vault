import ApplicationSerializer from '../application';

export default class KeymgmtProviderSerializer extends ApplicationSerializer {
  primaryKey = 'name';

  normalizeItems(payload) {
    let normalized = super.normalizeItems(payload);
    if (Array.isArray(normalized)) {
      return normalized.map((key) => ({
        id: key.name,
        name: key.name,
        backend: payload.backend,
      }));
    }
    return normalized;
  }

  serialize(snapshot) {
    const json = super.serialize(...arguments);
    return {
      ...json,
      credentials: snapshot.record.credentials,
    };
  }
}
