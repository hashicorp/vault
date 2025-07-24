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

const stanzaMaker = (path: string, policyStanzas: string[]) => {
  const caps = policyStanzas.length ? policyStanzas.map((c) => `"${c}"`).join(', ') : '';
  return `path "${path}" {
  capabilities = [${caps}]
}`;
};
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

class PolicyStanza {
  @tracked path: string;
  @tracked capabilities: string[] = [];

  constructor(path: string) {
    this.path = path;
  }

  get policyStanza() {
    return stanzaMaker(this.path, this.capabilities);
  }

  get hasCapabilities() {
    return this.capabilities.length !== 0;
  }

  @action
  setPermissions(event: HTMLElementEvent<HTMLInputElement>) {
    const { value, checked } = event.target;
    if (checked) {
      this.capabilities = addToArray(this.capabilities, value);
    } else {
      this.capabilities = removeFromArray(this.capabilities, value);
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

  existingPolicies: string[] | undefined = [];
  permissions = ['create', 'read', 'update', 'delete', 'list', 'patch', 'sudo'];
  identityOptions: Record<IdentitySelectionKey, Option[]> = {
    authMount: [],
    group: [],
    entity: [],
  };

  @tracked showFlyout = false;
  @tracked showPreview = false;
  @tracked policyAction = 'create';
  @tracked policyName = '';
  @tracked policyStanzas: PolicyStanza[] = [];

  @tracked selectedAssignments: Record<IdentitySelectionKey, Option[]> = {
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

  get assignmentText() {
    const identities = Object.keys(this.selectedAssignments)
      .filter((k) => this.selectedAssignments[k as IdentitySelectionKey].length)
      .map((type) => this.displayText(type).title.toLowerCase());

    if (identities.length > 1) {
      const lastItem = identities.pop();
      return `${identities.join(', ')} and ${lastItem}`;
    } else if (identities.length === 1) {
      return identities[0];
    } else {
      return 'identities';
    }
  }

  get policySnippet() {
    if (this.policyStanzas.length === 0) {
      return stanzaMaker('', []);
    }
    return this.formatPolicy(this.policyStanzas);
  }

  get cliSnippet() {
    return `vault policy write ${this.policyName || '[policy name]'} - <<EOF
${this.policySnippet}
EOF`;
  }

  get tfvpSnippet() {
    return `resource "vault_policy" "${this.policyName || '[policy name]'}" {
  name   = "${this.policyName || '[policy name]'}"
  policy = <<-EOT
${this.policySnippet}
EOT
}`;
  }

  @action
  openFlyout() {
    this.showFlyout = true;

    const { currentRoute, currentRouteName } = this.router;
    if (currentRoute && !currentRouteName?.includes('loading') && 'attributes' in currentRoute) {
      const { name, attributes } = currentRoute as { name: string; attributes: unknown };
      const apiPaths = mapApiPathToRoute(name);
      this.policyStanzas = apiPaths?.map((fn) => new PolicyStanza(fn(attributes))) || [];
      this.policyStanzas = [...this.policyStanzas];
    }
    return [];
  }

  @action
  handlePolicySelection(event: HTMLElementEvent<HTMLInputElement>) {
    const { name, value } = event.target;
    if (name === 'policyAction') {
      // either "create" or "edit"
      this.policyAction = value;
      // reset policy name
      this.policyName = '';
    } else {
      this.policyName = value;
    }
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
    const setOptions = (type: IdentitySelectionKey, options: Option[] | undefined) => {
      this.identityOptions[type] = options || [];
    };

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
    const item = new PolicyStanza('');
    this.policyStanzas.push(item);
    // Trigger an update
    this.policyStanzas = [...this.policyStanzas];
  }

  @action
  deletePath(path: string) {
    this.policyStanzas = [...this.policyStanzas.filter((c) => c.path !== path)];
  }

  // the magic part!
  @action
  async applyPolicy() {
    await this.createOrEditPolicy();
    // TODO
    // request to actually apply policies to identities
  }

  async createOrEditPolicy() {
    // policySnippet is purely for rendering the policy example. when it comes time to actually create/edit
    // the policy we want to remove any stanzas without permissions
    let policyPayload = this.formatPolicy(this.policyStanzas.filter((c) => c.hasCapabilities));
    // if editing an existing policy, fetch the original
    if (this.policyAction === 'edit') {
      const { policy, rules } = await this.api.sys.policiesReadAclPolicy2(this.policyName);
      // supposedly "rules" is deprecated, but that was the only key that returned data for me ¯\_(ツ)_/¯
      const data = policy || rules || '';
      // add existing policy to payload
      policyPayload = data.concat(`\n\n`, policyPayload);
    }
    try {
      await this.api.sys.policiesWriteAclPolicy2(this.policyName, { policy: policyPayload });
    } catch (error) {
      const { message } = await this.api.parseError(error);
      console.debug(message); // eslint-disable-line
    }
  }

  // HELPERS
  formatPolicy(policyStanzas: PolicyStanza[]) {
    return policyStanzas.map((c) => c.policyStanza).join('\n');
  }
}
