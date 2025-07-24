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
import mapApiPathToRoute from 'vault/utils/policy-path-map';
import { EntityListByNameListEnum, GroupListByNameListEnum } from '@hashicorp/vault-client-typescript';

import type ApiService from 'vault/services/api';
import type { HTMLElementEvent } from 'vault/forms';
import type RouterService from '@ember/routing/router-service';

interface Option {
  type: string;
  name: string;
  authType?: string;
}

const IDENTITY_TYPES = {
  authMount: 'Authentication mount',
  group: 'Group',
  entity: 'Entity',
} as const;

type IdentitySelectionKey = keyof typeof IDENTITY_TYPES;
// type IdentityOptionKey = (typeof IDENTITY_TYPES)[IdentitySelectionKey];
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

  @tracked selectedAssignments: Record<IdentitySelectionKey, Option[]> = {
    authMount: [],
    group: [],
    entity: [],
  };

  permissions = ['create', 'read', 'update', 'delete', 'list', 'patch', 'sudo'];

  identityOptions: Record<IdentitySelectionKey, Option[]> = {
    authMount: [],
    group: [],
    entity: [],
  };

  displayText = (type: string) => {
    switch (type) {
      case 'authMount':
        return {
          title: 'Authentication mounts',
          description: 'Policy will be applied to users who authenticate with the selected mounts.',
        };
      case 'group':
        return {
          title: 'Groups',
          description: 'Policy will be applied to users who belong to the selected groups.',
        };
      case 'entity':
        return {
          title: 'Entities',
          description: 'Policy will be applied to users who belong to the selected entities.',
        };

      default:
        return {
          title: `Select a ${type}`,
          description: 'The policy will be applied the selected resource.',
        };
    }
  };

  constructor(owner: unknown, args: Record<string, never>) {
    super(owner, args);
    this.fetchPolicies();
    this.fetchIdentities();
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

  @action
  selectPolicy(event: HTMLElementEvent<HTMLInputElement>) {
    const { value } = event.target;
    this.policyAction = value;
  }

  @action
  handleAssignment(type: IdentitySelectionKey, selection: Option[]) {
    this.selectedAssignments[type] = selection;
    // trigger DOM update
    this.selectedAssignments = Object.assign(this.selectedAssignments);
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
  async fetchIdentities() {
    const setOptions = (type: IdentitySelectionKey, options: Option[] | undefined) =>
      (this.identityOptions[type] = options || []);

    let type: IdentitySelectionKey;
    try {
      type = 'entity';
      const { keys } = await this.api.identity.entityListByName(EntityListByNameListEnum.TRUE);
      const entities = keys?.map((k) => ({ type, name: k }));
      setOptions(type, entities);
    } catch {
      // nope
    }

    try {
      type = 'group';
      const { keys } = await this.api.identity.groupListByName(GroupListByNameListEnum.TRUE);
      const groups = keys?.map((k) => ({ type, name: k }));
      setOptions(type, groups);
    } catch {
      // nope
    }

    try {
      type = 'authMount';
      const { auth } = await this.api.sys.internalUiListEnabledVisibleMounts();
      const mounts = this.api
        .responseObjectToArray(auth, 'path')
        .map((m) => ({ type, name: m.path, authType: m.type }));
      setOptions(type, mounts);
    } catch {
      // nope
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
