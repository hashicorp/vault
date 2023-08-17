import { FormField } from 'vault/app-types';
import PkiCertificateBaseModel from './base';

export default class PkiCertificateGenerateModel extends PkiCertificateBaseModel {
  name: string;
  formFields: FormField[];
  formFieldsGroup: {
    [k: string]: FormField[];
  }[];
}
