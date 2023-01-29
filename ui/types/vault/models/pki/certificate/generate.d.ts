import Model from '@ember-data/model';
import { FormField } from 'vault/app-types';

export default class PkiCertificateGenerateModel extends Model {
  name: string;
  backend: string;
  serialNumber: string;
  certificate: string;
  formFields: FormField[];
  formFieldsGroup: {
    [k: string]: FormField[];
  }[];
}
