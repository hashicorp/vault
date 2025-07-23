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

import type { HTMLElementEvent } from 'vault/forms';
import type RouterService from '@ember/routing/router-service';
import mapApiPathToRoute from 'vault/utils/policy-path-map';

class Capability {
  @tracked permissions: string[] = [];
  @tracked path: string;

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

  setPermissions(checked: boolean, value: string) {
    if (checked) {
      this.permissions = addToArray(this.permissions, value);
    } else {
      this.permissions = removeFromArray(this.permissions, value);
    }
  }
}

export default class PolicyBuilder extends Component {
  @service declare readonly router: RouterService;

  @tracked showFlyout = false;
  @tracked capabilities: Capability[] = [];

  permissions = ['create', 'read', 'update', 'delete', 'list', 'patch', 'sudo'];

  get policySnippet() {
    if (this.capabilities.length === 0) {
      return `path " " {
  permissions = [ ]
}`;
    }
    return this.capabilities.map((c) => c.policyStanza).join('\n');
  }

  @action
  updatePermissions(event: HTMLElementEvent<HTMLInputElement>) {
    const { name, value, checked } = event.target;
    let capability = this.capabilities.find((c) => c.path === name);
    if (!capability) {
      capability = new Capability(name);
      this.capabilities.push(capability);
    }
    capability.setPermissions(checked, value);
    // Trigger reactivity by reassigning the array
    // Remove any stanzas with no permissions
    this.capabilities = [...this.capabilities.filter((c) => c.hasPermissions)];
  }

  get paths() {
    const { currentRoute, currentRouteName } = this.router;
    if (currentRoute && !currentRouteName?.includes('loading') && 'attributes' in currentRoute) {
      const { name, attributes } = currentRoute as { name: string; attributes: unknown };
      const paths = mapApiPathToRoute(name);
      return paths?.map((fn) => fn(attributes)) || [];
    }
    return [];
  }

  get context() {
    return this.paths.join(', ');
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
