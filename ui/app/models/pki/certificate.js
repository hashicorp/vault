import Model, { attr } from '@ember-data/model';

export default class PkiCertificateModel extends Model {
  @attr('string', { readOnly: true }) backend;
  @attr('string') commonName;
  @attr('string') issueDate;
  @attr('string') serialNumber;
  @attr('string') notAfter;
  @attr('string') notBeforeDuration;
}
