import { encodePath } from 'vault/utils/path-encoding-helpers';
import ApplicationAdapter from '../../application';

export default class PkiCertificateGenerateAdapter extends ApplicationAdapter {
  namespace = 'v1';

  deleteRecord(store, type, snapshot) {
    const { backend, serialNumber, certificate } = snapshot.record;
    // Revoke certificate requires either serial_number or certificate
    const data = serialNumber ? { serial_number: serialNumber } : { certificate };
    return this.ajax(`${this.buildURL()}/${backend}/revoke`, 'POST', { data });
  }

  urlForCreateRecord(modelName, snapshot) {
    const { name, backend } = snapshot.record;
    if (!name || !backend) {
      throw new Error('URL for create record is missing required attributes');
    }
    return `${this.buildURL()}/${encodePath(backend)}/issue/${encodePath(name)}`;
  }
}
