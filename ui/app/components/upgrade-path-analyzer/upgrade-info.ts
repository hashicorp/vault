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

// TODO: improve typing
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
}

interface UpgradeInfoArgs {
  breadcrumbs: Array<Breadcrumb>;
  upgradeInfo?: UpgradeInfo[];
  tab?: string;
}

const PAGE_SIZE = 10;

export default class UpgradeInfoComponent extends Component<UpgradeInfoArgs> {
  @service declare readonly router: RouterService;

  @tracked currentPage = 1;

  get selectedTabIndex() {
    const index = parseInt(this.args.tab ?? '', 10);
    return !isNaN(index) && index >= 0 && index < this.tabs.length ? index : 0;
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

  @action
  paginatedItems(panelIndex: number) {
    const panel = this.panels[panelIndex] ?? [];
    return paginate(panel, { page: this.currentPage, pageSize: PAGE_SIZE });
  }

  get panels() {
    const knownIssues = this.upgradeInfo?.flatMap((item) => item.known_issues ?? []) ?? [];
    const breakingChanges = this.upgradeInfo?.flatMap((item) => item.breaking_changes ?? []) ?? [];
    const newBehavior = this.upgradeInfo?.flatMap((item) => item.new_behavior ?? []) ?? [];

    const formattedPanels = [
      this.formatKnownIssues(knownIssues),
      this.formatChangeOrBehavior(breakingChanges),
      this.formatChangeOrBehavior(newBehavior),
    ];

    return formattedPanels;
  }

  get tabs() {
    return [
      { text: 'Known issues', icon: 'shield-alert', count: this.panels[0]?.length },
      { text: 'Breaking changes', icon: 'alert-triangle', count: this.panels[1]?.length },
      { text: 'New behavior', icon: 'alert-circle', count: this.panels[2]?.length },
      { text: 'Rollback steps', icon: 'rewind' },
    ];
  }

  get upgradeInfo() {
    return this.args.upgradeInfo;
  }
}
