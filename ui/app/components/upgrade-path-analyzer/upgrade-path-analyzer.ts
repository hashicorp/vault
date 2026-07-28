/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { compareVersions } from 'vault/utils/version-utils';

import type ApiService from 'vault/services/api';
import type ClusterModel from 'vault/models/cluster';
import type CurrentClusterService from 'vault/services/current-cluster';
import type { HTMLElementEvent } from 'vault/forms';
import type RouterService from '@ember/routing/router-service';
import type { UpgradeInfo } from './upgrade-info';
import type VersionService from 'vault/services/version';

interface ReplicationSecondary {
  node_id?: string;
}

interface ReplicationInfo {
  clusterName: string;
  clusterId: string;
  drMode: string;
  drReplicationEnabled: boolean;
  drSecondaries: ReplicationSecondary[];
  perfMode: string;
  perfReplicationEnabled: boolean;
  perfSecondaries: ReplicationSecondary[];
}

enum Scenarios {
  SINGLE = 'single-instance',
  ENTERPRISE_REPLICATION = 'enterprise-replication',
}

interface UpgradePathAnalyzerArgs {
  onSetUpgradeInfo: (info: DisplayedInfo[]) => void;
}

type DisplayedInfo = Omit<UpgradeInfo, 'version'>;

export default class UpgradePathAnalyzer extends Component<UpgradePathAnalyzerArgs> {
  @service declare readonly api: ApiService;
  @service declare readonly router: RouterService;
  @service declare readonly currentCluster: CurrentClusterService;
  @service declare readonly version: VersionService;

  @tracked selectedVersion: string | null = null;
  @tracked upgradeInfo: DisplayedInfo[] | null = null;
  @tracked generalUpgradeInfoResp: UpgradeInfo[] | null = null;
  @tracked isLoading = false;

  constructor(owner: unknown, args: UpgradePathAnalyzerArgs) {
    super(owner, args);
    this.fetchReleaseInfo();
  }

  async fetchReleaseInfo() {
    try {
      // TODO: will update to this.api.sys.releaseInfoRead() once redundant suffix is addressed
      const { versions } = await this.api.sys.releaseInfoReadRead();
      if (versions) {
        this.generalUpgradeInfoResp = versions as UpgradeInfo[];
      }
    } catch (e) {
      console.error('Failed to fetch release info:', e);
    }
  }

  get targetVersions() {
    return ['1.19.0', '1.19.1', '1.19.5', '1.20.0', '1.20.1', '1.20.6', '1.21.7', '2.0.1'];
  }
  get currentVersion() {
    return this.version.version as unknown as string;
  }

  get breakingChangesCount(): number {
    const breakingChanges = this.upgradeInfo?.flatMap((item) => item.breaking_changes ?? []) ?? [];
    return breakingChanges.length;
  }

  get issueCount(): number {
    const knownIssues = this.upgradeInfo?.flatMap((item) => item.known_issues ?? []) ?? [];
    return knownIssues.length;
  }

  get newBehaviorCount(): number {
    const newBehavior = this.upgradeInfo?.flatMap((item) => item.new_behavior ?? []) ?? [];
    return newBehavior.length;
  }

  get cluster() {
    return this.currentCluster.cluster as ClusterModel | null;
  }

  /**
   * Provide normalized replication data used to generate order and checklist steps.
   */
  get clusterReplicationInfo(): ReplicationInfo | null {
    if (!this.cluster) {
      return null;
    }

    const drSecondaries = this.cluster.dr?.knownSecondaries ?? [];
    const perfSecondaries = this.cluster.performance?.knownSecondaries ?? [];

    return {
      clusterName: this.cluster.name ?? 'Current cluster',
      clusterId: this.cluster.dr?.clusterIdDisplay ?? this.cluster.performance?.clusterIdDisplay ?? 'unknown',
      drMode: this.cluster.dr?.mode ?? 'unknown',
      drReplicationEnabled: Boolean(this.cluster.dr?.replicationEnabled),
      drSecondaries,
      perfMode: this.cluster.performance?.mode ?? 'unknown',
      perfReplicationEnabled: Boolean(this.cluster.performance?.replicationEnabled),
      perfSecondaries,
    };
  }

  get scenario(): Scenarios {
    const info = this.clusterReplicationInfo;
    if (info?.drReplicationEnabled || info?.perfReplicationEnabled) {
      return Scenarios.ENTERPRISE_REPLICATION;
    }
    return Scenarios.SINGLE;
  }

  get upgradeGuidanceMessage(): string {
    return this.scenario === Scenarios.ENTERPRISE_REPLICATION
      ? 'General order: always upgrade secondary instances first, then primary instances.'
      : 'Single instance: upgrade the current Vault instance after creating a backup';
  }

  get upgradeOrder(): string[] {
    return this.scenario === Scenarios.ENTERPRISE_REPLICATION
      ? this.replicationUpgradeOrder
      : ['Upgrade the single Vault instance.'];
  }

  get upgradeSteps(): string[] {
    if (this.scenario === Scenarios.ENTERPRISE_REPLICATION) {
      return this.replicationUpgradeSteps;
    } else {
      return [
        'Create backup of Primary cluster via command vault operator raft snapshot save primary.snap on that cluster',
        'Stop Vault on the current instance.',
        `Install Vault ${this.selectedVersion} over the existing instance.`,
        'Start Vault.',
        'Unseal Vault if required.',
      ];
    }
  }

  get upgradeStepsDownloadText(): string {
    const orderLines = this.upgradeOrder.map((step, index) => `${index + 1}. ${step}`);
    const stepLines = this.upgradeSteps.map((step, index) => `${index + 1}. ${step}`);

    return [
      '# Vault upgrade steps',
      '',
      this.upgradeGuidanceMessage,
      '',
      '## Upgrade order',
      ...orderLines,
      '',
      '## Detailed upgrade steps',
      ...stepLines,
    ].join('\n');
  }

  get replicationUpgradeOrder(): string[] {
    const info = this.clusterReplicationInfo;
    if (!info) {
      return ['Upgrade secondary clusters first, then primary clusters.'];
    }

    const order: string[] = [];

    info.drSecondaries.forEach((secondary) => {
      order.push(`Upgrade DR secondary: ${secondary.node_id}`);
    });

    info.perfSecondaries.forEach((secondary) => {
      order.push(`Upgrade performance secondary: ${secondary.node_id}`);
    });

    if (this.isPrimaryMode(info.drMode) || this.isPrimaryMode(info.perfMode)) {
      order.push(`Upgrade primary cluster: ${info.clusterName} (${info.clusterId}).`);
    }

    if (!order.length) {
      order.push('Upgrade secondary clusters first, then primary clusters.');
    }

    return order;
  }

  // Steps for replicated deployments
  get replicationUpgradeSteps(): string[] {
    const info = this.clusterReplicationInfo;
    if (!info) {
      return [
        `Upgrade secondary clusters first, then upgrade Vault ${this.selectedVersion} on the primary cluster.`,
      ];
    }

    const steps: string[] = [];

    info.drSecondaries.forEach((secondary) => {
      steps.push(
        `Create backup of ${secondary.node_id} (DR Secondary) via command vault operator raft snapshot save ${secondary.node_id}.snap on that cluster`
      );
    });

    info.perfSecondaries.forEach((secondary) => {
      steps.push(
        `Create backup of ${secondary.node_id} (Perf Secondary) via command vault operator raft snapshot save ${secondary.node_id}.snap on that cluster`
      );
    });

    if (this.isPrimaryMode(info.drMode) || this.isPrimaryMode(info.perfMode)) {
      steps.push(
        `Create backup of Primary cluster via command vault operator raft snapshot save primary.snap on that cluster`
      );
    }

    steps.push(`Back up the current Vault configuration`);
    steps.push(`Perform any prerequisites noted in the documentation`);
    steps.push(`Use SIGINT or SIGTERM to shut down Vault`);
    steps.push(`Install ${this.selectedVersion}`);
    steps.push('Start Vault');

    return steps;
  }

  @action
  async onAnalyzeClick() {
    if (this.generalUpgradeInfoResp && this.selectedVersion) {
      this.isLoading = true;
      await new Promise((resolve) => setTimeout(resolve, 600));
      this.upgradeInfo = this.filterReleaseInfo(
        this.generalUpgradeInfoResp,
        this.currentVersion,
        this.selectedVersion
      );
      this.args.onSetUpgradeInfo(this.upgradeInfo);
      this.isLoading = false;
    }
  }

  @action
  onVersionSelect(event: HTMLElementEvent<HTMLInputElement>) {
    const { value } = event.target;
    this.selectedVersion = value;
    this.upgradeInfo = null;
  }

  /**
   * Filter the release info response according to the current and selected versions.
   *
   * - Filters breaking_changes to only include those introduced between currentVersion
   *   (exclusive) and targetVersion (inclusive).
   * - Filters new_behavior to only include those introduced between currentVersion
   *   (exclusive) and targetVersion (inclusive).
   * - Filters known_issues to only include those:
   *   - Found between currentVersion (exclusive) and targetVersion (inclusive)
   *   - AND not yet fixed by targetVersion (fixed === 'No' OR fixed version > targetVersion)
   *
   * Returns a single aggregated entry so that the existing flatMap consumers work unchanged.
   */
  filterReleaseInfo(releaseInfo: UpgradeInfo[], current: string, target: string): DisplayedInfo[] {
    const inRange = (v: string) => compareVersions(v, current) > 0 && compareVersions(v, target) <= 0;
    const notFixedByTarget = (fixed: string) => fixed === 'No' || compareVersions(fixed, target) > 0;

    const breaking_changes = releaseInfo
      .flatMap((e) => e.breaking_changes ?? [])
      .filter((item) => inRange(item.introduced));

    const new_behavior = releaseInfo
      .flatMap((e) => e.new_behavior ?? [])
      .filter((item) => inRange(item.introduced));

    const known_issues = releaseInfo
      .flatMap((e) => e.known_issues ?? [])
      .filter((item) => inRange(item.found) && notFixedByTarget(item.fixed));

    return [{ breaking_changes, new_behavior, known_issues }];
  }

  private isPrimaryMode(mode: string): boolean {
    return mode === 'primary';
  }
}
