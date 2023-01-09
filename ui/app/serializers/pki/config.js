import ApplicationSerializer from '../application';

export default class PkiConfigSerializer extends ApplicationSerializer {
  primaryKey = 'request_id';
  attrs = {
    formType: { serialize: false },
  };

  /*
  serializing: ID for import could be `imported_issuers[0]`?
  */
}
