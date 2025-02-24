import { FormField, WithFormFieldsAndValidationsModel } from 'vault/app-types';

export default interface SshCaConfig extends WithFormFieldsAndValidationsModel {
  backend: string;
  privateKey: string;
  publicKey: string;
  generateSigningKey: boolean;

  configurableParams: ['privateKey', 'publicKey', 'generateSigningKey'];

  get displayAttrs(): FormField[];
}
