/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

import type { UnauthMountsByType } from 'vault/vault/auth/form';
import type { HTMLElementEvent } from 'vault/forms';

interface Args {
  authTabData: UnauthMountsByType;
  handleTabClick: CallableFunction;
  selectedAuthMethod: string;
}

export default class AuthTabs extends Component<Args> {
  @tracked selectedMountPath = '';

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    this.setSelectedMountPath();
  }

  get mountDescription() {
    const selectedAuth = this.args.selectedAuthMethod;
    const mounts = this.args.authTabData[selectedAuth];
    const mount = mounts?.find((m) => m.path === this.selectedMountPath);
    return mount?.description ?? '';
  }

  get tabTypes() {
    return this.args.authTabData ? Object.keys(this.args.authTabData) : [];
  }

  get selectedTabIndex() {
    const index = this.tabTypes.indexOf(this.args.selectedAuthMethod);
    // negative index means the selected method isn't a tab, default to first tab
    return index < 0 ? 0 : index;
  }

  @action
  onClickTab(_event: HTMLElementEvent<HTMLInputElement>, idx: number) {
    const newMethod = this.tabTypes[idx];
    this.args.handleTabClick(newMethod);
    // Reset selected mount path when tab changes
    this.setSelectedMountPath();
  }

  @action
  setMount(event: HTMLElementEvent<HTMLInputElement>) {
    this.selectedMountPath = event.target.value;
  }

  private setSelectedMountPath() {
    const mounts = this.args.authTabData[this.args.selectedAuthMethod];
    const firstMount = mounts?.length ? mounts[0] : null;
    this.selectedMountPath = firstMount?.path ?? '';
  }
}
