import Model from '@ember-data/model';
import PkiCertificateGenerateModel from 'vault/models/pki/certificate/generate';
import PkiConfigImportModel from 'vault/models/pki/config/import';

declare module 'ember-data/types/registries/model' {
  export default interface ModelRegistry {
    'pki/config/import': PkiConfigImportModel;
    'pki/certificate/generate': PkiCertificateGenerateModel;
    // Catchall for any other models
    [key: string]: any;
  }
}
