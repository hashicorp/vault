import ApplicationSerializer from '../application';

export default class OidcScopeSerializer extends ApplicationSerializer {
  primaryKey = 'name';
}
