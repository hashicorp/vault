/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { service } from '@ember/service';
import type NamespaceService from 'vault/services/namespace';

interface Project {
  name: string;
  error?: string;
}

interface Org {
  name: string;
  projects: Project[];
  error?: string;
}

class Block {
  @tracked global = '';
  @tracked orgs: Org[] = [{ name: '', projects: [{ name: '' }] }];
  @tracked globalError = '';

  constructor(global = '', orgs: Org[] = [{ name: '', projects: [{ name: '' }] }]) {
    this.global = global;
    this.orgs = orgs;
  }

  validateInput(value: string): string {
    if (value.includes('/')) {
      return '"/" is not allowed in namespace names';
    } else if (value.includes(' ')) {
      return 'spaces are not allowed in namespace names';
    }
    return '';
  }
}

interface Args {
  wizardState: {
    namespacePaths: string[] | null;
    namespaceBlocks: Block[] | null;
  };
  updateWizardState: (key: string, value: unknown) => void;
}

export default class WizardNamespacesStepTemp extends Component<Args> {
  @service declare namespace: NamespaceService;
  @tracked blocks: Block[];
  duplicateErrorMessage = 'No duplicate namespaces names are allowed within the same level';

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    this.blocks = args.wizardState.namespaceBlocks || [new Block()];
  }

  get treeChartOptions() {
    const currentNamespace = this.namespace.currentNamespace || 'root';
    return {
      height: '400px',
      tree: {
        type: 'tree',
        rootTitle: currentNamespace,
      },
    };
  }

  get hasErrors(): boolean {
    return this.blocks.some((block) => {
      // Check valid nesting
      if (!this.isValidNesting(block)) return true;
      // Check global error
      if (block.globalError) return true;
      // Check org errors
      if (block.orgs.some((org) => org.error)) return true;
      // Check project errors
      return block.orgs.some((org) => org.projects.some((project) => project.error));
    });
  }

  isValidNesting(block: Block) {
    // If there are non-empty orgs but no global, then it is invalid
    if (block.orgs.some((org) => org.name.trim()) && !block.global.trim()) {
      return false;
    }

    // Check all projects have proper parents (global and org)
    return block.orgs.every((org) => {
      const hasProjects = org.projects.some((project) => project.name.trim());
      return !hasProjects || (block.global.trim() && org.name.trim());
    });
  }

  checkForDuplicateGlobals() {
    const globals = this.blocks.map((block) => block.global.trim()).filter((global) => global !== '');
    const globalCounts = new Map();

    globals.forEach((global) => {
      globalCounts.set(global, (globalCounts.get(global) || 0) + 1);
    });

    this.blocks.forEach((block) => {
      if (!block.globalError && globalCounts.get(block.global) > 1) {
        block.globalError = this.duplicateErrorMessage;
      } else if (globalCounts.get(block.global) === 1 && block.globalError === this.duplicateErrorMessage) {
        // remove outdated error message
        block.globalError = '';
      }
    });
  }

  updateWizardState() {
    this.args.updateWizardState('namespacePaths', this.hasErrors ? null : this.namespacePaths);
    this.args.updateWizardState('namespaceBlocks', this.hasErrors ? null : this.blocks);
  }

  @action
  addBlock() {
    this.blocks = [...this.blocks, new Block()];
  }

  @action
  deleteBlock(index: number) {
    if (this.blocks.length > 1) {
      this.blocks = this.blocks.filter((_, i) => i !== index);
    } else {
      // Reset the only remaining block to initial state
      this.blocks = [new Block()];
    }
    // Re-validate duplicate globals in case a duplicate was deleted
    this.checkForDuplicateGlobals();
    this.updateWizardState();
  }

  @action
  updateGlobalValue(blockIndex: number, event: Event) {
    const target = event.target as HTMLInputElement;
    const block = this.blocks[blockIndex];
    if (block) {
      block.global = target.value;
      block.globalError = block.validateInput(target.value);
      this.checkForDuplicateGlobals();
      this.updateWizardState();
    }
  }

  @action
  updateOrgValue(block: Block, orgToUpdate: Org, event: Event) {
    const target = event.target as HTMLInputElement;
    const value = target.value;
    const isDuplicate = block.orgs.some((org) => org !== orgToUpdate && org.name === value);

    const updatedOrgs = block.orgs.map((org) => {
      if (org === orgToUpdate) {
        return {
          ...org,
          name: value,
          error: isDuplicate ? this.duplicateErrorMessage : block.validateInput(value),
        };
      }
      return org;
    });
    block.orgs = updatedOrgs;

    // Trigger tree reactivity by reassigning the blocks array
    this.blocks = [...this.blocks];
    this.updateWizardState();
  }

  @action
  addOrg(block: Block) {
    block.orgs = [...block.orgs, { name: '', projects: [{ name: '' }] }];
  }

  @action
  removeOrg(block: Block, orgToRemove: Org) {
    if (block.orgs.length <= 1) return;
    block.orgs = block.orgs.filter((org) => org !== orgToRemove);

    // Trigger tree reactivity
    this.blocks = [...this.blocks];
  }

  @action
  updateProjectValue(block: Block, org: Org, projectToUpdate: Project, event: Event) {
    const target = event.target as HTMLInputElement;
    const value = target.value;
    const isDuplicate = org.projects.some((project) => project !== projectToUpdate && project.name === value);

    const updatedOrgs = block.orgs.map((currentOrg) => {
      if (currentOrg === org) {
        return {
          ...currentOrg,
          projects: currentOrg.projects.map((project) => {
            if (project === projectToUpdate) {
              return {
                name: value,
                error: isDuplicate ? this.duplicateErrorMessage : block.validateInput(value),
              };
            }
            return project;
          }),
        };
      }
      return currentOrg;
    });
    block.orgs = updatedOrgs;

    // Trigger tree reactivity by reassigning the blocks array
    this.blocks = [...this.blocks];
    this.updateWizardState();
  }

  @action
  addProject(block: Block, org: Org) {
    const updatedOrgs = block.orgs.map((currentOrg) => {
      if (currentOrg === org) {
        return {
          ...currentOrg,
          projects: [...currentOrg.projects, { name: '' }],
        };
      }
      return currentOrg;
    });
    block.orgs = updatedOrgs;
  }

  @action
  removeProject(block: Block, org: Org, projectToRemove: Project) {
    if (org.projects.length <= 1) return;

    const updatedOrgs = block.orgs.map((currentOrg) => {
      if (currentOrg === org) {
        return {
          ...currentOrg,
          projects: currentOrg.projects.filter((project) => project !== projectToRemove),
        };
      }
      return currentOrg;
    });
    block.orgs = updatedOrgs;

    // Trigger tree reactivity
    this.blocks = [...this.blocks];
  }

  get treeData() {
    const parsed = this.blocks.map((block) => {
      return {
        name: block.global,
        children: block.orgs
          .filter((org) => org.name.trim() !== '')
          .map((org) => {
            return {
              name: org.name,
              children: org.projects
                .filter((project) => project.name.trim() !== '')
                .map((project) => {
                  return {
                    name: project.name,
                  };
                }),
            };
          }),
      };
    });

    return parsed;
  }

  // The Carbon tree chart only supports displaying nodes with at least 1 "fork" i.e. at least 2 globals, 2 orgs or 2 projects
  get shouldShowTreeChart(): boolean {
    // Count total globals across blocks
    const globalsCount = this.blocks.filter((block) => block.global.trim() !== '').length;

    // Check if there are multiple globals
    if (globalsCount > 1) {
      return true;
    }

    // Check if any block has multiple orgs
    const hasMultipleOrgs = this.blocks.some(
      (block) => block.orgs.filter((org) => org.name.trim() !== '').length > 1
    );

    if (hasMultipleOrgs) {
      return true;
    }

    // Check if any org has multiple projects
    const hasMultipleProjects = this.blocks.some((block) =>
      block.orgs.some((org) => org.projects.filter((project) => project.name.trim() !== '').length > 1)
    );

    return hasMultipleProjects;
  }

  // Store namespace paths to be used for code snippets in the format "global", "global/org", "global/org/project"
  get namespacePaths(): string[] {
    return this.blocks
      .map((block) => {
        const results: string[] = [];

        // Add global namespace if it exists
        if (block.global.trim() !== '') {
          results.push(block.global);
        }

        block.orgs.forEach((org) => {
          if (org.name.trim() !== '') {
            // Add global/org namespace
            const globalOrg = [block.global, org.name].filter((value) => value.trim() !== '').join('/');
            if (globalOrg && !results.includes(globalOrg)) {
              results.push(globalOrg);
            }

            org.projects.forEach((project) => {
              if (project.name.trim() !== '') {
                // Add global/org/project namespace
                const fullNamespace = [block.global, org.name, project.name]
                  .filter((value) => value.trim() !== '')
                  .join('/');
                if (fullNamespace && !results.includes(fullNamespace)) {
                  results.push(fullNamespace);
                }
              }
            });
          }
        });
        return results;
      })
      .flat()
      .filter((namespace) => namespace !== '');
  }
}
