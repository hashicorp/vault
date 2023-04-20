import Model from '@ember-data/model';
import { FormField } from 'vault/app-types';

export default class PkiCrlModel extends Model {
  autoRebuildData: object;
  deltaCrlBuildingData: object;
  crlExpiryData: {
    enabled: boolean;
    duration: string;
  };
  ocspExpiryData: object;
  formFields: FormField[];
  get canSet(): boolean;
}
