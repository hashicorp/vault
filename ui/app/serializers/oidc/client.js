import ApplicationSerializer from '../application';

export default class OidcClientSerializer extends ApplicationSerializer {
  primaryKey = 'name';
}
