/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */
import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';
import type RouterService from '@ember/routing/router-service';
import type { Breadcrumb } from 'vault/vault/app-types';
import type { HTMLElementEvent } from 'vault/forms';
import { paginate } from 'core/utils/paginate-list';
import { formatDownloadText } from './upgrade-utils';

interface KnownIssue {
  found: string;
  fixed: string;
  workaround: string;
  edition: string;
  issue: string;
  link: string;
}

interface BreakingChange {
  edition: string;
  recommendations: boolean;
  introduced: string;
  change: string;
  link: string;
}

interface NewBehavior {
  edition: string;
  recommendations: boolean;
  introduced: string;
  change: string;
  link: string;
}

export interface UpgradeInfo {
  version: string;
  known_issues?: KnownIssue[];
  breaking_changes?: BreakingChange[];
  new_behavior?: NewBehavior[];
  rollback_steps?: string[];
  rollbackOrder?: string[];
  rollbackGuidanceMessage?: string;
}

interface UpgradeInfoArgs {
  breadcrumbs: Array<Breadcrumb>;
  upgradeInfo?: UpgradeInfo;
  tab?: string;
  targetVersion?: string;
}

const PAGE_SIZE = 10;

export default class UpgradeInfoComponent extends Component<UpgradeInfoArgs> {
  @service declare readonly router: RouterService;

  @tracked currentPage = 1;

  get panels() {
    const formattedPanels = [
      this.formatKnownIssues(this.upgradeInfo?.known_issues ?? []),
      this.formatChangeOrBehavior(this.upgradeInfo?.breaking_changes ?? []),
      this.formatChangeOrBehavior(this.upgradeInfo?.new_behavior ?? []),
    ];

    return formattedPanels;
  }

  get fullDownloadText() {
    let text = `${this.rollbackStepsDownloadText}\n`;

    this.upgradeInfo?.known_issues?.forEach((issue) => {
      text += `\nKnown issue: ${issue.issue}. Found in ${issue.found}. Fixed: ${issue.fixed}. Edition ${issue.edition}. Workaround: ${issue.workaround}. Link: ${issue.link}\n`;
    });

    this.upgradeInfo?.breaking_changes?.forEach((change) => {
      text += `\nBreaking change: ${change.change}. Introduced in ${change.introduced}. Edition: ${change.edition}. Recommendations: ${change.recommendations}.  Link: ${change.link}\n`;
    });

    this.upgradeInfo?.new_behavior?.forEach((behavior) => {
      text += `\nNew behavior: ${behavior.change}. Introduced in ${behavior.introduced}. Edition: ${behavior.edition}. Recommendations: ${behavior.recommendations}.  Link: ${behavior.link}\n`;
    });

    return text;
  }

  get rollbackStepsDownloadText(): string {
    const text = formatDownloadText(
      this.upgradeInfo?.rollbackOrder ?? [],
      this.upgradeInfo?.rollback_steps ?? [],
      this.upgradeInfo?.rollbackGuidanceMessage ?? '',
      'Rollback'
    );

    return text;
  }

  get selectedTabIndex() {
    const index = parseInt(this.args.tab ?? '', 10);
    return !isNaN(index) && index >= 0 && index < this.tabs.length ? index : 0;
  }

  get tabs() {
    return [
      { label: 'Known issues', icon: 'shield-alert', count: this.panels[0]?.length },
      { label: 'Breaking changes', icon: 'alert-triangle', count: this.panels[1]?.length },
      { label: 'New behavior', icon: 'alert-circle', count: this.panels[2]?.length },
      {
        label: 'Rollback steps',
        icon: 'rewind',
        count: this.args.upgradeInfo?.rollback_steps?.length ?? '0',
      },
    ];
  }

  get upgradeInfo() {
    return this.args.upgradeInfo;
  }

  @action
  paginatedItems(panelIndex: number) {
    const panel = this.panels[panelIndex] ?? [];
    return paginate(panel, { page: this.currentPage, pageSize: PAGE_SIZE });
  }

  @action
  onClickTab(_event: HTMLElementEvent<HTMLInputElement>, index: number) {
    this.currentPage = 1;
    this.router.replaceWith({ queryParams: { tab: String(index) } });
  }

  @action
  onPageChange(page: number) {
    this.currentPage = page;
  }

  private formatChangeOrBehavior(items: BreakingChange[] | NewBehavior[]) {
    return items.map((data) => ({
      badges: [
        { text: data.edition, color: 'neutral' },
        {
          text: data.recommendations ? 'Recommendations available' : 'No recommendation',
          color: data.recommendations ? 'highlight' : 'critical',
        },
      ],
      title: data.change,
      description: `Introduced in ${data.introduced}`,
      link: data.link,
    }));
  }

  private formatKnownIssues(issues: KnownIssue[]) {
    return issues.map((data) => {
      const isFixed = data.fixed !== 'No';
      const hasWorkaround = data.workaround === 'Yes';

      return {
        badges: [
          { text: data.edition, color: 'neutral' },
          {
            text: isFixed ? `Fixed in ${data.fixed}` : 'Not fixed',
            color: isFixed ? 'success' : 'critical',
          },
          {
            text: hasWorkaround ? 'Workaround available' : 'No workaround available',
            color: hasWorkaround ? 'neutral' : 'critical',
          },
        ],
        title: data.issue,
        description: `Found in ${data.found}`,
        link: data.link,
      };
    });
  }
}
