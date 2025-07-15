/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
import Component from '@glimmer/component';
import './plugin-card.scss';
import type { RegistryPlugin } from '../../../types/vault-registry';
export interface Plugin {
    Args: {
        plugin: RegistryPlugin;
    };
    Blocks: {
        default: [];
    };
    Element: HTMLElement;
}
export default class PluginCard extends Component<Plugin> {
    get isEnterprisePlugin(): boolean;
    get pluginPublishDate(): string;
}
//# sourceMappingURL=plugin-card.d.ts.map