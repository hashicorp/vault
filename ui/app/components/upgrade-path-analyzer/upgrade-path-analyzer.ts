/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { formatDownloadText } from './upgrade-utils';
import { cleanVersion, compareVersions, parseVersion } from 'vault/utils/version-utils';

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
  onSetUpgradeInfo: (info: DisplayedInfo) => void;
}

export type DisplayedInfo = Omit<UpgradeInfo, 'version'> & {
  rollback_steps: string[];
  rollbackOrder: string[];
  rollbackGuidanceMessage: string;
  targetVersion: string;
};

export default class UpgradePathAnalyzer extends Component<UpgradePathAnalyzerArgs> {
  @service declare readonly api: ApiService;
  @service declare readonly router: RouterService;
  @service declare readonly currentCluster: CurrentClusterService;
  @service declare readonly version: VersionService;

  @tracked selectedVersion: string | null = null;
  @tracked upgradeInfo: DisplayedInfo | null = null;
  @tracked generalUpgradeInfoResp: UpgradeInfo[] | null = null;
  @tracked isLoading = false;
  @tracked isModalOpen = false;
  @tracked hasError = false;
  @tracked targetVersions: string[] = [];

  constructor(owner: unknown, args: UpgradePathAnalyzerArgs) {
    super(owner, args);
    this.fetchReleaseInfo();
  }

  async fetchReleaseInfo() {
    try {
      const { versions } = await this.api.sys.releaseInfoReadReleaseInfo();
      if (versions) {
        this.generalUpgradeInfoResp = versions as UpgradeInfo[];
      }
    } catch (e) {
      this.hasError = true;
    }

    try {
      const { versions } = await this.api.sys.vaultVersionsReadRead('enterprise');
      this.targetVersions = versions as unknown as string[];
    } catch (e) {
      this.hasError = true;
    }
  }

  get cards() {
    return [
      {
        icon: 'shield-alert',
        title: 'Known issues',
        description: 'These are all the known issues documented with the version selected.',
        count: this.upgradeInfo?.known_issues?.length ?? 0,
      },
      {
        icon: 'alert-triangle',
        title: 'Breaking changes',
        description: 'These are functional changes from one version to the other.',
        count: this.upgradeInfo?.breaking_changes?.length ?? 0,
      },
      {
        icon: 'gift',
        title: 'New behavior',
        description: 'New behavior introduced and released in the version selected.',
        count: this.upgradeInfo?.new_behavior?.length ?? 0,
      },
      {
        icon: 'rewind',
        title: 'Rollback steps',
        description: 'Follow these steps to safely rollback.',
        count: this.upgradeInfo?.rollback_steps.length ?? 0,
      },
    ];
  }

  /**
   * From the full fetched list, keep only:
   * 1. Versions strictly greater than the current version.
   * 2. For each MAJOR.MINOR group, only the highest PATCH (the latest release
   *    of that minor line).
   */
  get filteredTargetVersions(): string[] {
    const current = this.currentVersion;
    if (!this.targetVersions.length) return [];
    // If the current version is not yet known, show nothing — the getter will
    // recompute once version.version is populated (it is @tracked on the service).
    if (!current) return [];

    // Step 1: drop anything <= current version
    const newer = this.targetVersions.filter((v) => compareVersions(v, current) > 0);

    // Step 2: keep only the highest patch per MAJOR.MINOR group
    const latestPerMinor = new Map<string, string>();
    for (const v of newer) {
      const [major = 0, minor = 0] = parseVersion(v);
      const key = `${major}.${minor}`;
      const existing = latestPerMinor.get(key);
      if (!existing || compareVersions(v, existing) > 0) {
        latestPerMinor.set(key, v);
      }
    }

    // Return in ascending order (Map preserves insertion order after sort)
    return Array.from(latestPerMinor.values()).sort((a, b) => compareVersions(a, b));
  }

  get currentVersion() {
    const raw = this.version.version as string | null;
    return raw ? cleanVersion(raw) : '';
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

  get rollbackGuidanceMessage(): string {
    return this.scenario === Scenarios.ENTERPRISE_REPLICATION
      ? 'General order: always rollback secondary instances first, then primary instances.'
      : 'Single instance: rollback the current Vault instance from your backup';
  }

  get rollbackOrder(): string[] {
    return this.scenario === Scenarios.ENTERPRISE_REPLICATION
      ? this.replicationRollbackOrder
      : ['Rollback the single Vault instance.'];
  }

  get rollbackSteps(): string[] {
    return this.scenario === Scenarios.ENTERPRISE_REPLICATION
      ? this.replicationRollbackSteps
      : [
          'Stop the Vault service via command sudo systemctl stop vault or via command Stop-Process <vault_pid>',
          `Install your previous version of Vault (${this.currentVersion}) over your existing instance.`,
          'Replace the upgraded Vault data store with your pre-upgrade snapshot.',
          'Replace the upgraded Vault configuration with your pre-upgrade configuration.',
          'Start Vault.',
          'Verify the current version via command vault status | grep Version',
          'Unseal vault',
          'Test the rollback',
        ];
  }

  get upgradeStepsDownloadText(): string {
    const text = formatDownloadText(
      this.upgradeOrder,
      this.upgradeSteps,
      this.upgradeGuidanceMessage,
      'Upgrade'
    );

    return text;
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

  get replicationRollbackOrder(): string[] {
    const info = this.clusterReplicationInfo;
    if (!info) {
      return ['Rollback secondary clusters first, then primary clusters.'];
    }

    const order: string[] = [];

    info.drSecondaries.forEach((secondary) => {
      order.push(`Rollback DR secondary: ${secondary.node_id}`);
    });

    info.perfSecondaries.forEach((secondary) => {
      order.push(`Rollback performance secondary: ${secondary.node_id}`);
    });

    if (this.isPrimaryMode(info.drMode) || this.isPrimaryMode(info.perfMode)) {
      order.push(`Rollback primary cluster: ${info.clusterName} (${info.clusterId}).`);
    }

    if (!order.length) {
      order.push('Rollback secondary clusters first, then primary clusters.');
    }

    return order;
  }

  get replicationRollbackSteps(): string[] {
    const info = this.clusterReplicationInfo;
    if (!info) {
      return [
        `Rollback secondary clusters first, then rollback to Vault ${this.selectedVersion} on the primary cluster.`,
      ];
    }

    const steps: string[] = [];

    info.drSecondaries.forEach((secondary) => {
      steps.push(`Rollback from backup of ${secondary.node_id} (DR Secondary).`);
    });

    info.perfSecondaries.forEach((secondary) => {
      steps.push(`Rollback from backup of ${secondary.node_id} (Perf Secondary).`);
    });

    if (this.isPrimaryMode(info.drMode) || this.isPrimaryMode(info.perfMode)) {
      steps.push('Rollback from backup of Primary cluster.');
    }

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
  filterReleaseInfo(
    releaseInfo: UpgradeInfo[],
    current: string,
    target: string,
    targetVersion = target
  ): DisplayedInfo {
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

    return {
      breaking_changes,
      new_behavior,
      known_issues,
      rollback_steps: this.rollbackSteps,
      rollbackOrder: this.rollbackOrder,
      rollbackGuidanceMessage: this.rollbackGuidanceMessage,
      targetVersion,
    };
  }
  private isPrimaryMode(mode: string): boolean {
    return mode === 'primary';
  }
}
