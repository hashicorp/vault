import Model, { attr, hasMany } from '@ember-data/model';
import ArrayProxy from '@ember/array/proxy';
import PromiseProxyMixin from '@ember/object/promise-proxy-mixin';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import { withModelValidations } from 'vault/decorators/model-validations';
import { isPresent } from '@ember/utils';

const validations = {
  name: [
    { type: 'presence', message: 'Name is required.' },
    // ARG TODO add in after Claire pushes her branch
    // {
    //   type: 'containsWhiteSpace',
    //   message: 'Name cannot contain whitespace.',
    // },
  ],
  targets: [
    {
      validator(model) {
        const entityIds = model.hasMany('entity_ids').ids();
        const groupIds = model.hasMany('group_ids').ids();
        return isPresent(entityIds) || isPresent(groupIds);
      },
      message: 'At least one entity or group is required.',
    },
  ],
};

@withModelValidations(validations)
export default class OidcAssignmentModel extends Model {
  @attr('string') name;
  @hasMany('identity/entity') entity_ids;
  @hasMany('identity/group') group_ids;

  @lazyCapabilities(apiPath`identity/oidc/assignment/${'name'}`, 'name') assignmentPath;
  @lazyCapabilities(apiPath`identity/oidc/assignment`) assignmentsPath;

  get targets() {
    return ArrayProxy.extend(PromiseProxyMixin).create({
      promise: this.prepareTargets(),
    });
  }

  async prepareTargets() {
    const targets = [];

    for (const key of ['entity_ids', 'group_ids']) {
      (await this[key]).forEach((model) => {
        targets.addObject({
          key,
          title: model.name,
          subTitle: model.id,
        });
      });
    }

    return targets;
  }

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
