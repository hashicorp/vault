import ApplicationSerializer from '../application';

export default class PkiConfigSerializer extends ApplicationSerializer {
  primaryKey = 'request_id';
  attrs = {
    formType: { serialize: false },
  };
}
