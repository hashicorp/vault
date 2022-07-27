import Model, { attr } from '@ember-data/model';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import { withModelValidations } from 'vault/decorators/model-validations';
import { isPresent } from '@ember/utils';

const validations = {
  name: [
    { type: 'presence', message: 'Name is required.' },
    // {
    //   type: 'containsWhiteSpace',
    //   message: 'Name cannot contain whitespace.',
    // },
  ],
  targets: [
    {
      validator(model) {
        return isPresent(model.entityIds) || isPresent(model.groupIds);
      },
      message: 'At least one entity or group is required.',
    },
  ],
};

@withModelValidations(validations)
export default class OidcAssignmentModel extends Model {
  @attr('string') name;
  @attr('array') entityIds;
  @attr('array') groupIds;

  // CAPABILITIES
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
