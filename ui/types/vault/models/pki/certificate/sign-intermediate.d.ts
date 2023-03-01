import PkiCertificateBaseModel from './base';
import { FormField, FormFieldGroups, ModelValidations } from 'vault/app-types';
export default class PkiCertificateSignIntermediateModel extends PkiCertificateBaseModel {
  role: string;
  csr: string;
  formFields: FormField[];
  formFieldGroups: FormFieldGroups;
  issuerRef: string;
  maxPathLength: string;
  notBeforeDuration: string;
  permittedDnsDomains: string;
  useCsrValues: boolean;
  usePss: boolean;
  skid: string;
  signatureBits: string;
  validate(): ModelValidations;
}
