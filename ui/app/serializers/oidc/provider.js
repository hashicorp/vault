import ApplicationSerializer from '../application';

export default class OidcProviderSerializer extends ApplicationSerializer {
  primaryKey = 'name';
}
