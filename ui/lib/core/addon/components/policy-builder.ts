/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';
import { addToArray } from 'vault/helpers/add-to-array';
import { removeFromArray } from 'vault/helpers/remove-from-array';

import type ApiService from 'vault/services/api';
import type { HTMLElementEvent } from 'vault/forms';
import type RouterService from '@ember/routing/router-service';
import mapApiPathToRoute from 'vault/utils/policy-path-map';

class Capability {
  @tracked path: string;
  @tracked permissions: string[] = [];

  constructor(path: string) {
    this.path = path;
  }

  get policyStanza() {
    return `path "${this.path}" {
  permissions = [ ${this.permissions.map((c) => `"${c}"`).join(', ')} ]
}`;
  }

  get hasPermissions() {
    return this.permissions.length !== 0;
  }

  @action
  setPermissions(event: HTMLElementEvent<HTMLInputElement>) {
    const { value, checked } = event.target;
    if (checked) {
      this.permissions = addToArray(this.permissions, value);
    } else {
      this.permissions = removeFromArray(this.permissions, value);
    }
  }

  @action
  setPath(event: HTMLElementEvent<HTMLInputElement>) {
    this.path = event.target.value;
  }
}

export default class PolicyBuilder extends Component {
  @service declare readonly api: ApiService;
  @service declare readonly router: RouterService;

  @tracked showFlyout = false;
  @tracked policyAction = 'create';
  @tracked policyName = '';
  @tracked existingPolicies: string[] | undefined = [];
  @tracked capabilities: Capability[] = [];
  @tracked showAdvanced = false;

  @tracked identities: string[] | undefined = [];

  permissions = ['create', 'read', 'update', 'delete', 'list', 'patch', 'sudo'];
  // identityTypes = [
  //   { type: 'authMount', label: 'Authentication mount' },
  //   { type: 'group', label: 'Group' },
  //   { type: 'entity', label: 'Entity' },
  // ];

  identityOptions = {
    'Authentication mount': [{ id: 'userpass/' }, { id: 'oidc/' }, { id: 'ldap/' }],
    Group: [{ id: 'admins' }, { id: 'platform-engineers' }, { id: 'sales' }],
    Entity: [{ id: 'bob' }, { id: 'matilda' }, { id: 'lorraine' }],
  };

  constructor(owner: unknown, args: Record<string, never>) {
    super(owner, args);
    this.fetchPolicies();
  }

  get context() {
    const params = this.router.currentRoute?.parent?.params;
    return params ? Object.values(params).join('/') : '';
  }

  get policySnippet() {
    if (this.capabilities.length === 0) {
      return `path " " {
  permissions = [ ]
}`;
    }
    return this.capabilities.map((c) => c.policyStanza).join('\n');
  }

  @action
  openFlyout() {
    this.showFlyout = true;

    const { currentRoute, currentRouteName } = this.router;
    if (currentRoute && !currentRouteName?.includes('loading') && 'attributes' in currentRoute) {
      const { name, attributes } = currentRoute as { name: string; attributes: unknown };
      const apiPaths = mapApiPathToRoute(name);
      this.capabilities = apiPaths?.map((fn) => new Capability(fn(attributes))) || [];
      this.capabilities = [...this.capabilities];
    }
    return [];
  }

  // @action
  // handleAssignment(event: HTMLElementEvent<HTMLInputElement>) {
  //   // do something
  // }

  @action
  selectPolicy(event: HTMLElementEvent<HTMLInputElement>) {
    const { value } = event.target;
    this.policyAction = value;
  }

  @action
  async fetchPolicies() {
    try {
      const { keys } = await this.api.sys.policiesListAclPolicies2();
      this.existingPolicies = keys;
    } catch {
      // nah
    }
  }

  @action
  addPath() {
    const item = new Capability('');
    this.capabilities.push(item);
    // Trigger an update
    this.capabilities = [...this.capabilities];
  }

  @action
  deletePath(path: string) {
    this.capabilities = [...this.capabilities.filter((c) => c.path !== path)];
  }

  tfvp = `
resource "vault_auth_backend" "userpass" {
  type = "userpass"
}

resource "vault_generic_endpoint" "danielle-user" {
   path                 = "auth/path/users/danielle-vault-user"
   ignore_absent_fields = true
   data_json = <<EOT
{
   "token_policies": ["developer-vault-policy"],
   "password": "Vividness Itinerary Mumbo Reassure"
}
EOT
}
`;
  cli = `-
vault policy write top-secret-policy - <<EOF
path "sys/auth" {
  capabilities = ["read"]
}
EOF
path "/v2/admin/secret/data/top-secret" {
  capabilities = [ "update" ]
}
`;
}
