import ApplicationSerializer from '../application';

export default class OidcKeySerializer extends ApplicationSerializer {
  primaryKey = 'name';
}
