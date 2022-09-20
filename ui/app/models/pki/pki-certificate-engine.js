import Model, { attr } from '@ember-data/model';

export default class PkiCertificatesEngineModel extends Model {
  @attr('string', { readOnly: true }) backend;
  @attr('string') commonName;
  @attr('string') issueDate;
  @attr('string') serialNumber;
  @attr('string') notAfter;
  @attr('string') notBeforeDuration;

  // ARG TODO these are the READ/details view of the certificate. Will need to add to for the create.
}
