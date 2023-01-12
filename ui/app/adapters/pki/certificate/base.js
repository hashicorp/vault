import { encodePath } from 'vault/utils/path-encoding-helpers';
import ApplicationAdapter from '../../application';

export default class PkiCertificateBaseAdapter extends ApplicationAdapter {
  namespace = 'v1';

  deleteRecord(store, type, snapshot) {
    const { backend, serialNumber, certificate } = snapshot.record;
    // Revoke certificate requires either serial_number or certificate
    const data = serialNumber ? { serial_number: serialNumber } : { certificate };
    return this.ajax(`${this.buildURL()}/${encodePath(backend)}/revoke`, 'POST', { data });
  }
}
