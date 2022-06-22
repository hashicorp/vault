import Model, { attr, hasMany } from '@ember-data/model';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';

export default class OidcAssignmentModel extends Model {
  @attr('string') name;
  @hasMany('identity/entity') entityIds;
  @hasMany('identity/group') groupIds;

  @lazyCapabilities(apiPath`identity/oidc/assignment/${'name'}`, 'name') assignmentPath;
  @lazyCapabilities(apiPath`identity/oidc/assignment`) assignmentsPath;

  get canCreate() {
    return this.assignmentPath.get('canCreate');
  }
  get canRead() {
    return this.assignmentPath.get('canRead');
  }
  get canEdit() {
    return this.assignmentPath.get('canUpdate');
  }
  get canDelete() {
    return this.assignmentPath.get('canDelete');
  }
  get canList() {
    return this.assignmentsPath.get('canList');
  }

  @lazyCapabilities(apiPath`identity/entity`) entitiesPath;
  get canListEntities() {
    return this.entitiesPath.get('canList');
  }

  @lazyCapabilities(apiPath`identity/group`) groupsPath;
  get canListGroups() {
    return this.groupsPath.get('canList');
  }
}
