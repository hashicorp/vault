import { computed } from '@ember/object';
import DS from 'ember-data';

import AuthConfig from '../auth-config';
import fieldToAttrs from 'vault/utils/field-to-attrs';
import { combineFieldGroups } from 'vault/utils/openapi-to-attrs';

const { attr } = DS;

export default AuthConfig.extend({
  useOpenAPI: true,
  binddn: attr('string', {
    helpText: 'Used when performing user search. Example: cn=vault,ou=Users,dc=example,dc=com',
  }),
  bindpass: attr('string', {
    helpText: 'Used along with binddn when performing user search',
    sensitive: true,
  }),
  userdn: attr('string', {
    helpText: 'Base DN under which to perform user search. Example: ou=Users,dc=example,dc=com',
  }),
  userattr: attr('string', {
    helpText:
      'Attribute on user attribute object matching the username passed when authenticating. Examples: sAMAccountName, cn, uid',
  }),
  upndomain: attr('string', {
    helpText:
      'The userPrincipalDomain used to construct the UPN string for the authenticating user. The constructed UPN will appear as [username]@UPNDomain. Example: example.com, which will cause vault to bind as username@example.com.',
  }),

  groupfilter: attr('string', {
    helpText:
      'Go template used when constructing the group membership query. The template can access the following context variables: [UserDN, Username]. The default is (|(memberUid={{.Username}})(member={{.UserDN}})(uniqueMember={{.UserDN}})), which is compatible with several common directory schemas. To support nested group resolution for Active Directory, instead use the following query: (&(objectClass=group)(member:1.2.840.113556.1.4.1941:={{.UserDN}}))',
  }),
  groupdn: attr('string', {
    helpText:
      'LDAP search base for group membership search. This can be the root containing either groups or users. Example: ou=Groups,dc=example,dc=com',
  }),
  groupattr: attr('string', {
    helpText:
      'LDAP attribute to follow on objects returned by groupfilter in order to enumerate user group membership. Examples: for groupfilter queries returning group objects, use: cn. For queries returning user objects, use: memberOf. The default is cn.',
  }),
  useTokenGroups: attr('boolean', {
    helpText:
      'Use the Active Directory tokenGroups constructed attribute to find the group memberships. This returns all security groups for the user, including nested groups. In an Active Directory environment with a large number of groups this method offers increased performance. Selecting this will cause Group DN, Attribute, and Filter to be ignored.',
  }),

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
