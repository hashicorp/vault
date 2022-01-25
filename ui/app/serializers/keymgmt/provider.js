import ApplicationSerializer from '../application';

export default class KeymgmtProviderSerializer extends ApplicationSerializer {
  primaryKey = 'name';

  serialize(snapshot) {
    const json = super.serialize(...arguments);
    return {
      ...json,
      credentials: snapshot.record.credentials,
    };
  }
}
