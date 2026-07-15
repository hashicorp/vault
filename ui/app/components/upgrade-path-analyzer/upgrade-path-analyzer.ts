/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';

//TODO: Move data later once endpoint is available
//Note: data can't be pulled from tests directory
import { UPGRADE_INFO } from 'vault/constants/upgrade-info';

import type { HTMLElementEvent } from 'vault/forms';
import type VersionService from 'vault/services/version';
import type RouterService from '@ember/routing/router-service';

interface UpgradePathAnalyzerArgs {
  onSetUpgradeInfo: (info: unknown[]) => void;
}

export default class UpgradePathAnalyzer extends Component<UpgradePathAnalyzerArgs> {
  @service declare readonly version: VersionService;
  @service declare readonly router: RouterService;
  @tracked selectedVersion: string | null = null;
  @tracked upgradeInfo: unknown[] | null = null;
  get targetVersions() {
    return ['1.20.0', '1.20.1'];
  }
  get currentVersion() {
    return this.version.version;
  }

  get issueCount(): number {
    const knownIssues = this.upgradeInfo?.flatMap((item: any) => item.known_issues ?? []) ?? [];
    return knownIssues.length;
  }

  get breakingChangesCount(): number {
    const breakingChanges = this.upgradeInfo?.flatMap((item: any) => item.breaking_changes ?? []) ?? [];
    return breakingChanges.length;
  }

  get newBehaviorCount(): number {
    const newBehavior = this.upgradeInfo?.flatMap((item: any) => item.new_behavior ?? []) ?? [];
    return newBehavior.length;
  }

  @action
  onVersionSelect(event: HTMLElementEvent<HTMLInputElement>) {
    const { value } = event.target;
    this.selectedVersion = value;
  }

  @action
  onAnalyzeClick() {
    this.upgradeInfo = UPGRADE_INFO;
    this.args.onSetUpgradeInfo(this.upgradeInfo);
  }
}
