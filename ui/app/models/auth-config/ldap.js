import { computed } from '@ember/object';
import DS from 'ember-data';

import AuthConfig from '../auth-config';
import fieldToAttrs from 'vault/utils/field-to-attrs';

const { attr } = DS;

export default AuthConfig.extend({
  url: attr('string', {
    label: 'URL',
  }),
  starttls: attr('boolean', {
    defaultValue: false,
    label: 'Issue StartTLS command after establishing an unencrypted connection',
  }),
  tlsMinVersion: attr('string', {
    label: 'Minimum TLS Version',
    defaultValue: 'tls12',
    possibleValues: ['tls10', 'tls11', 'tls12'],
  }),

  tlsMaxVersion: attr('string', {
    label: 'Maximum TLS Version',
    defaultValue: 'tls12',
    possibleValues: ['tls10', 'tls11', 'tls12'],
  }),
  insecureTls: attr('boolean', {
    defaultValue: false,
    label: 'Skip LDAP server SSL certificate verification',
  }),
  certificate: attr('string', {
    label: 'CA certificate to verify LDAP server certificate',
    editType: 'file',
  }),

  binddn: attr('string', {
    label: 'Name of Object to bind (binddn)',
    helpText: 'Used when performing user search. Example: cn=vault,ou=Users,dc=example,dc=com',
  }),
  bindpass: attr('string', {
    label: 'Password',
    helpText: 'Used along with binddn when performing user search',
    sensitive: true,
  }),

  userdn: attr('string', {
    label: 'User DN',
    helpText: 'Base DN under which to perform user search. Example: ou=Users,dc=example,dc=com',
  }),
  userattr: attr('string', {
    label: 'User Attribute',
    defaultValue: 'cn',
    helpText:
      'Attribute on user attribute object matching the username passed when authenticating. Examples: sAMAccountName, cn, uid',
  }),
  discoverdn: attr('boolean', {
    defaultValue: false,
    label: 'Use anonymous bind to discover the bind DN of a user',
  }),
  denyNullBind: attr('boolean', {
    defaultValue: true,
    label: 'Prevent users from bypassing authentication when providing an empty password',
  }),
  upndomain: attr('string', {
    label: 'User Principal (UPN) Domain',
    helpText:
      'The userPrincipalDomain used to construct the UPN string for the authenticating user. The constructed UPN will appear as [username]@UPNDomain. Example: example.com, which will cause vault to bind as username@example.com.',
  }),

  groupfilter: attr('string', {
    label: 'Group Filter',
    helpText:
      'Go template used when constructing the group membership query. The template can access the following context variables: [UserDN, Username]. The default is (|(memberUid={{.Username}})(member={{.UserDN}})(uniqueMember={{.UserDN}})), which is compatible with several common directory schemas. To support nested group resolution for Active Directory, instead use the following query: (&(objectClass=group)(member:1.2.840.113556.1.4.1941:={{.UserDN}}))',
  }),
  groupdn: attr('string', {
    label: 'Group DN',
    helpText:
      'LDAP search base for group membership search. This can be the root containing either groups or users. Example: ou=Groups,dc=example,dc=com',
  }),
  groupattr: attr('string', {
    label: 'Group Attribute',
    defaultValue: 'cn',

    helpText:
      'LDAP attribute to follow on objects returned by groupfilter in order to enumerate user group membership. Examples: for groupfilter queries returning group objects, use: cn. For queries returning user objects, use: memberOf. The default is cn.',
  }),
  useTokenGroups: attr('boolean', {
    defaultValue: false,
    label: 'Use Token Groups',
    helpText:
      'Use the Active Directory tokenGroups constructed attribute to find the group memberships. This returns all security groups for the user, including nested groups. In an Active Directory environment with a large number of groups this method offers increased performance. Selecting this will cause Group DN, Attribute, and Filter to be ignored.',
  }),

  fieldGroups: computed(function() {
    const groups = [
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
    return fieldToAttrs(this, groups);
  }),
});
