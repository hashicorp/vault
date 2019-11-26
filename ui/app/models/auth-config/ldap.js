import { computed } from '@ember/object';

import AuthConfig from '../auth-config';
import fieldToAttrs from 'vault/utils/field-to-attrs';
import { combineFieldGroups } from 'vault/utils/openapi-to-attrs';

export default AuthConfig.extend({
  useOpenAPI: true,
  fieldGroups: computed(function() {
    let groups = [
      {
        default: ['url'],
      },
      {
        'LDAP Options': [
          'starttls',
          'insecureTls',
          'discoverdn',
          'denyNullBind',
          'tlsMinVersion',
          'tlsMaxVersion',
          'certificate',
          'userattr',
          'upndomain',
        ],
      },
      {
        'Customize User Search': ['binddn', 'userdn', 'bindpass'],
      },
      {
        'Customize Group Membership Search': ['groupfilter', 'groupattr', 'groupdn', 'useTokenGroups'],
      },
    ];
    if (this.newFields) {
      groups = combineFieldGroups(groups, this.newFields, []);
    }
    return fieldToAttrs(this, groups);
  }),
});
