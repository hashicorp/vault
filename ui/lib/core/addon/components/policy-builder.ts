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
import type RouterService from '@ember/routing/router-service';
import type FlashMessages from 'vault/services/flash-messages';
import type { HTMLElementEvent } from 'vault/forms';

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
  // authMount: 'Authentication mount',
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
  @service declare readonly flashMessages: FlashMessages;
  @service declare readonly router: RouterService;

  existingPolicies: string[] | undefined = [];
  permissions = ['create', 'read', 'update', 'delete', 'list', 'patch', 'sudo'];
  identityOptions: Record<IdentitySelectionKey, Option[]> = {
    // authMount: [],
    group: [],
    entity: [],
  };

  @tracked error = '';
  @tracked showFlyout = false;
  @tracked showPreview = false;
  @tracked policyAction = 'create';
  @tracked policyName = '';
  @tracked existingPolicy = ''; // if a policy is being edited
  @tracked policyStanzas: PolicyStanza[] = [];
  @tracked selectedAssignments: Record<IdentitySelectionKey, Option[]> = {
    // authMount: [],
    group: [],
    entity: [],
  };

  displayText = (type: string) => {
    switch (type) {
      // case 'authMount':
      //   return {
      //     title: 'Authentication mounts',
      //     description: 'Policy will apply to users who authenticate with the selected mounts.',
      //   };
      case 'group':
        return {
          title: 'Groups',
          description: 'Policy will apply to users who belong to the selected groups.',
        };
      case 'entity':
        return {
          title: 'Entities',
          description: 'Policy will apply to users who belong to the selected entities.',
        };
      default:
        return {
          title: `Select a ${type}`,
          description: 'The policy will apply the selected resource.',
        };
    }
  };

  constructor(owner: unknown, args: Record<string, never>) {
    super(owner, args);
    this.fetchPolicies();
    this.fetchIdentities();
  }

  get applySubtext() {
    const identities = this.filteredAssignments.map((type) => this.displayText(type).title.toLowerCase());

    if (identities.length > 1) {
      const lastItem = identities.pop();
      return `${identities.join(', ')} and ${lastItem}`;
    } else if (identities.length === 1) {
      return identities[0];
    } else {
      return 'identities';
    }
  }

  get filteredAssignments() {
    return Object.keys(this.selectedAssignments).filter(
      (k) => this.selectedAssignments[k as IdentitySelectionKey].length
    );
  }

  get policySnippet() {
    if (this.policyStanzas.length === 0) {
      return stanzaMaker('', []);
    }
    return this.formatPolicy(this.policyStanzas);
  }

  get actualPolicy() {
    // policySnippet is purely for rendering the preview. when it comes time to use the actually
    // policy we want to remove any stanzas without permissions
    let actualPolicy = this.formatPolicy(this.policyStanzas.filter((c) => c.hasCapabilities));
    // if editing an existing policy, add the policy data
    if (this.policyAction === 'edit') {
      const lineBreak = actualPolicy ? `\n\n` : '';
      actualPolicy = this.existingPolicy.concat(lineBreak, actualPolicy);
    }
    return actualPolicy;
  }

  get cliSnippet() {
    if (this.policyName) {
      const cliCommand = (o: Option, type: IdentitySelectionKey) =>
        `vault write identity/${type} name="${o.name}" policies="default, ${this.policyName}"`;
      const command = this.buildAssignmentSnippet(cliCommand);

      return `vault policy write ${this.policyName} - <<EOF
${this.actualPolicy}
EOF
${command}`;
    } else {
      return '# Select a policy or fill in a name to preview commands!';
    }
  }

  get tfvpSnippet() {
    if (this.policyName) {
      const tfvpCommand = (o: Option, type: IdentitySelectionKey) => {
        return `resource "vault_identity_${type}" "${o.name}" {
 name     = "${o.name}"
 policies = ["default", "${this.policyName}"]
}`;
      };

      const command = this.buildAssignmentSnippet(tfvpCommand);
      return `resource "vault_policy" "${this.policyName}" {
  name   = "${this.policyName}"
  policy = <<-EOT
${this.actualPolicy}
EOT
}
${command}`;
    } else {
      return '# Select a policy or fill in a name to preview commands!';
    }
  }

  @action
  openFlyout() {
    this.showFlyout = true;

    const { currentRoute, currentRouteName } = this.router;
    if (currentRoute && !currentRouteName?.includes('loading') && 'attributes' in currentRoute) {
      const { name, attributes } = currentRoute as { name: string; attributes: unknown };
      const apiPaths = mapApiPathToRoute(name);
      // hardcoding the check for "backend" since this is hackweek and
      // only secrets are supported
      // try the parent if there are none at the current route
      let params;
      if (
        attributes &&
        typeof attributes === 'object' &&
        !Array.isArray(attributes) &&
        'backend' in attributes
      ) {
        params = attributes;
      } else {
        params = currentRoute?.parent?.params;
      }
      this.policyStanzas = apiPaths?.map((fn) => new PolicyStanza(fn(params))) || [];
      this.policyStanzas = [...this.policyStanzas];
    }
    return [];
  }

  @action
  async handleRadio(event: HTMLElementEvent<HTMLInputElement>) {
    // value is either "create" or "edit"
    this.policyAction = event.target.value;
    // reset policy name
    this.policyName = '';
  }

  @action
  async handleCreatePolicy(event: HTMLElementEvent<HTMLInputElement>) {
    this.policyName = event.target.value;
    // reset existing policy in case "edit" was previously selected
    this.existingPolicy = '';
  }

  @action
  async handleEditPolicy(name: string) {
    this.policyName = name;
    const { policy, rules } = await this.api.sys.policiesReadAclPolicy2(this.policyName);
    // supposedly "rules" is deprecated, but that was the only key that returned data for me ¯\_(ツ)_/¯
    this.existingPolicy = policy || rules || '';
  }

  @action
  handleAssignment(type: IdentitySelectionKey, selection: Option[]) {
    this.selectedAssignments[type] = selection || [];
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

    // try {
    //   type = 'authMount';
    //   const { auth } = await this.api.sys.internalUiListEnabledVisibleMounts();
    //   const mounts = this.api
    //     .responseObjectToArray(auth, 'path')
    //     .map((m) => ({ type, name: m.path, authType: m.type }));
    //   setOptions(type, mounts);
    // } catch {
    //   // nope
    // }
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
    // TODO: request to actually apply policies to groups/entities
  }

  async createOrEditPolicy() {
    try {
      const policyPayload = this.actualPolicy;
      await this.api.sys.policiesWriteAclPolicy2(this.policyName, { policy: policyPayload });
      const word = this.policyAction === 'create' ? 'created' : 'updated';
      this.flashMessages.success(`Success! The policy: ${this.policyName} has been successfully ${word}!`);
      this.resetState();
    } catch (error) {
      const { message } = await this.api.parseError(error);
      this.error = message;
      console.debug(message); // eslint-disable-line
    }
  }

  @action
  resetState() {
    this.showFlyout = false;
    this.showPreview = false;
    this.policyAction = 'create';
    this.policyName = '';
    this.selectedAssignments = {
      // authMount: [],
      group: [],
      entity: [],
    };
  }

  // HELPERS
  formatPolicy(policyStanzas: PolicyStanza[]) {
    return policyStanzas.map((c) => c.policyStanza).join('\n');
  }

  buildAssignmentSnippet(commandTemplate: CallableFunction) {
    let assignments: string[] = [];
    if (this.filteredAssignments.length) {
      for (const [key, value] of Object.entries(this.selectedAssignments)) {
        if (!value?.length) continue;
        if (key === 'authMount') {
          // do auth mount command
        } else {
          const commands = value.map((g) => commandTemplate(g, key as IdentitySelectionKey));
          assignments = [...commands, ...assignments];
        }
      }
      return assignments.length ? assignments.join('\n') : '';
    }
    return '';
  }
}
