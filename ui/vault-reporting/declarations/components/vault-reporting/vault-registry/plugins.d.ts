/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
import Component from '@glimmer/component';
import type { RegistryPlugin } from '../../../types/vault-registry';
export interface PluginList {
    Args: {
        plugins?: RegistryPlugin[];
    };
    Blocks: {
        default: [];
    };
    Element: HTMLElement;
}
export default class Plugins extends Component<PluginList> {
    get plugins(): RegistryPlugin[] | undefined;
}
//# sourceMappingURL=plugins.d.ts.map