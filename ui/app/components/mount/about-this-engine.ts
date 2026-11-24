/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { getAboutEngineInfo } from 'vault/utils/about-engine-info';
import type { AboutEngineEntry, BulletPart } from 'vault/utils/about-engine-info';

interface Args {
  /** The normalized engine type, e.g. 'kv', 'aws', 'azure', 'gcp', 'database' */
  engineType: string;
}

interface BulletPartResolved {
  isLink: boolean;
  isCode: boolean;
  text: string;
  href?: string;
}

export interface EngineInfoRow {
  key: string;
  label: string;
  icon: string;
  bgClass: string;
  labelColor: string;
  iconColor: string;
  bullets: BulletPartResolved[][];
}

/**
 * @module Mount::AboutThisEngine
 *
 * Collapsible "About this engine" card for the five common secrets engines (KV, AWS, Azure,
 * GCP, Database). Shows three rows: preferred use cases, key features, and limitations.
 * Open by default; renders nothing for engine types not in `about-engine-info`.
 *
 * @example
 * ```hbs
 * <Mount::AboutThisEngine @engineType="kv" />
 * ```
 */
export default class MountAboutThisEngineComponent extends Component<Args> {
  /** Controls whether the three info rows are visible. Open by default. */
  @tracked isOpen = true;

  get engineInfo(): AboutEngineEntry | null {
    return getAboutEngineInfo(this.args.engineType);
  }

  get hasEngineInfo(): boolean {
    return this.engineInfo !== null;
  }

  get toggleLabel(): string {
    return this.isOpen ? 'Hide details' : 'View details';
  }

  @action
  toggleOpen(): void {
    this.isOpen = !this.isOpen;
  }

  get rows(): EngineInfoRow[] {
    const info = this.engineInfo;
    if (!info) return [];

    return [
      {
        key: 'preferredFor',
        label: 'Preferred for',
        icon: 'check-circle',
        bgClass: 'has-background-surface-highlight',
        labelColor: 'highlight-on-surface',
        iconColor: 'highlight-on-surface',
        bullets: this.resolveBullets(info.preferredFor),
      },
      {
        key: 'keyFeatures',
        label: 'Key features',
        icon: 'star',
        bgClass: 'background-neutral-50',
        labelColor: 'primary',
        iconColor: 'primary',
        bullets: this.resolveBullets(info.keyFeatures),
      },
      {
        key: 'notSuitedFor',
        label: 'Not suited for',
        icon: 'alert-triangle',
        bgClass: 'has-background-surface-warning',
        labelColor: 'warning-on-surface',
        iconColor: 'warning-on-surface',
        bullets: this.resolveBullets(info.notSuitedFor),
      },
    ];
  }

  private resolveBullets(bullets: BulletPart[][]): BulletPartResolved[][] {
    return bullets.map((parts) =>
      parts.map((part): BulletPartResolved => {
        if (typeof part === 'object' && 'href' in part) {
          return { isLink: true, isCode: false, text: part.text, href: part.href };
        }
        if (typeof part === 'object' && 'code' in part) {
          return { isLink: false, isCode: true, text: part.code };
        }
        return { isLink: false, isCode: false, text: part };
      })
    );
  }
}
