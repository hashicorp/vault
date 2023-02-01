import PkiCertificateBaseModel from './certificate/base';
import { FormField, FormFieldGroups, ModelValidations } from 'vault/app-types';
export default class PkiIssuerModel extends PkiCertificateBaseModel {
  useOpenAPI(): boolean;
  issuerId: string;
  keyId: string;
  uriSans: string;
  issuerName: string;
  leafNotAfterBehavior: string;
  usage: string;
  manualChain: string;
  issuingCertificates: string;
  crlDistributionPoints: string;
  ocspServers: string;
  /** these are all instances of the capabilities model which should be converted to native class and typed
  rotateExported: any;
  rotateInternal: any;
  rotateExisting: any;
  crossSignPath: any;
  signIntermediate: any;
  -------------------- **/
  formFields: Array<FormField>;
  formFieldGroups: FormFieldGroups;
  allFields: Array<FormField>;
  get canRotateIssuer(): boolean;
  get canCrossSign(): boolean;
  get canSignIntermediate(): boolean;
  get canConfigure(): boolean;
  validate(): ModelValidations;
}
