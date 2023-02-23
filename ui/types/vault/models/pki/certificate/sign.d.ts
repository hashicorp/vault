import PkiCertificateBaseModel from './base';
import { FormField, FormFieldGroups, ModelValidations } from 'vault/app-types';
export default class PkiCertificateSignModel extends PkiCertificateBaseModel {
  name: string;
  formFields: FormField[];
  formFieldGroups: FormFieldGroups;
  validate(): ModelValidations;
}
