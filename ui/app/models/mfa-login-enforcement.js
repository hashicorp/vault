import Model, { attr, hasMany } from '@ember-data/model';
import ArrayProxy from '@ember/array/proxy';
import PromiseProxyMixin from '@ember/object/promise-proxy-mixin';
import { methods } from 'vault/helpers/mountable-auth-methods';
import { withModelValidations } from 'vault/decorators/model-validations';
import { isPresent } from '@ember/utils';

const validations = {
  name: [{ type: 'presence', message: 'Name is required' }],
  mfa_methods: [{ type: 'presence', message: 'At least one MFA method is required' }],
  targets: [
    {
      validator(model) {
        // avoid async fetch of records here and access relationship ids to check for presence
        const entityIds = model.hasMany('identity_entities').ids();
        const groupIds = model.hasMany('identity_groups').ids();
        return (
          isPresent(model.auth_method_accessors) ||
          isPresent(model.auth_method_types) ||
          isPresent(entityIds) ||
          isPresent(groupIds)
        );
      },
      message:
        "At least one target is required. If you've selected one, click 'Add' to make sure it's added to this enforcement.",
    },
  ],
};
@withModelValidations(validations)
export default class MfaLoginEnforcementModel extends Model {
  @attr('string') name;
  @hasMany('mfa-method') mfa_methods;
  @attr('string') namespace_id;
  @attr('array', { defaultValue: () => [] }) auth_method_accessors; // ["auth_approle_17a552c6"]
  @attr('array', { defaultValue: () => [] }) auth_method_types; // ["userpass"]
  @hasMany('identity/entity') identity_entities;
  @hasMany('identity/group') identity_groups;

  get targets() {
    return ArrayProxy.extend(PromiseProxyMixin).create({
      promise: this.prepareTargets(),
    });
  }

  async prepareTargets() {
    const mountableMethods = methods(); // use for icon lookup
    let authMethods;
    const targets = [];

    if (this.auth_method_accessors.length || this.auth_method_types.length) {
      // fetch all auth methods and lookup by accessor to get mount path and type
      try {
        const { data } = await this.store.adapterFor('auth-method').findAll();
        authMethods = Object.keys(data).map((key) => ({ path: key, ...data[key] }));
      } catch (error) {
        // swallow this error
      }
    }

    if (this.auth_method_accessors.length) {
      const selectedAuthMethods = authMethods.filter((model) => {
        return this.auth_method_accessors.includes(model.accessor);
      });
      targets.addObjects(
        selectedAuthMethods.map((method) => {
          const mount = mountableMethods.findBy('type', method.type);
          const icon = mount.glyph || mount.type;
          return {
            icon,
            link: 'vault.cluster.access.method',
            linkModels: [method.path.slice(0, -1)],
            title: method.path,
            subTitle: method.accessor,
          };
        })
      );
    }

    this.auth_method_types.forEach((type) => {
      const mount = mountableMethods.findBy('type', type);
      const icon = mount.glyph || mount.type;
      const mountCount = authMethods.filterBy('type', type).length;
      targets.addObject({
        key: 'auth_method_types',
        icon,
        title: type,
        subTitle: `All ${type} mounts (${mountCount})`,
      });
    });

    for (const key of ['identity_entities', 'identity_groups']) {
      (await this[key]).forEach((model) => {
        targets.addObject({
          key,
          icon: 'user',
          link: 'vault.cluster.access.identity.show',
          linkModels: [key.split('_')[1], model.id, 'details'],
          title: model.name,
          subTitle: model.id,
        });
      });
    }

    return targets;
  }
}
