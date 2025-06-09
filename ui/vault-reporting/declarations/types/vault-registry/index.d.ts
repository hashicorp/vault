/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
import type { HdsBadgeSignature } from '@hashicorp/design-system-components/components/hds/badge/index';
type PluginTag = 'built-in' | 'official' | 'partner' | 'community' | 'enterprise';
export interface RegistryPlugin {
    id: string;
    pluginName: string;
    author: string;
    description: string;
    externalUrl: string;
    pluginType: 'AUTH' | 'SECRET' | 'DATABASE';
    pluginVersion: string;
    isRegistered?: boolean;
    official: OfficialPlugin;
    tags: PluginTag;
    publishDate: Date;
}
interface OfficialPlugin {
    author: HdsBadgeSignature['Args']['icon'];
    tags: string;
}
export {};
//# sourceMappingURL=index.d.ts.map