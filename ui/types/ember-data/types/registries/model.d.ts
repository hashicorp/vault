import Model from '@ember-data/model';
import PkiActionModel from 'vault/models/pki/action';
import PkiCertificateGenerateModel from 'vault/models/pki/certificate/generate';

declare module 'ember-data/types/registries/model' {
  export default interface ModelRegistry {
    'pki/action': PkiActionModel;
    'pki/certificate/generate': PkiCertificateGenerateModel;
    // Catchall for any other models
    [key: string]: any;
  }
}
