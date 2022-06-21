import Model, { attr, hasMany } from '@ember-data/model';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';

export default class OidcAssignmentModel extends Model {
  @attr('string') name;
  @hasMany('identity/entity') identity_entities;
  @hasMany('identity/group') identity_groups;

  @attr('array', {
    label: 'Entities',
    editType: 'searchSelect',
    fallbackComponent: 'string-list',
    models: ['identity/entity'],
  })
  entityIds;

  @attr('array', {
    label: 'Groups',
    editType: 'searchSelect',
    fallbackComponent: 'string-list',
    models: ['identity/group'],
  })
  groupIds;

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
}
