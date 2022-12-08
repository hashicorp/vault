import { encodePath } from 'vault/utils/path-encoding-helpers';
import ApplicationAdapter from '../../application';

export default class PkiCertificateGenerateAdapter extends ApplicationAdapter {
  namespace = 'v1';

  createRecord(store, type, snapshot) {
    const roleName = snapshot.attr('name');
    const url = `${this.buildURL()}/${encodePath(snapshot.record.backend)}/issue/${roleName}`;
    const data = store.serializerFor(type.modelName).serialize(snapshot);
    return this.ajax(url, 'POST', {
      data,
    });
  }
}
