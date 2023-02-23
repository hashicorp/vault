import Model from '@ember-data/model';
import { FormField, FormFieldGroups, ModelValidations } from 'vault/app-types';
export default class PkiIssuerModel extends Model {
  useOpenAPI(): boolean;
  issuerId: string;
  issuerName: string;
  issuerRef(): string;
  keyId: string;
  uriSans: string;
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
  formFields: FormField[];
  formFieldGroups: FormFieldGroups;
  allFields: FormField[];
  get canRotateIssuer(): boolean;
  get canCrossSign(): boolean;
  get canSignIntermediate(): boolean;
  get canConfigure(): boolean;
  validate(): ModelValidations;
}
