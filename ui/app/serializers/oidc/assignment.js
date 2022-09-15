import ApplicationSerializer from '../application';

export default class OidcAssignmentSerializer extends ApplicationSerializer {
  primaryKey = 'name';
}
