/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
type ISODateString = `${number}${number}-${number}${number}-${number}${number}`;
type ISOTimeString = `${number}${number}:${number}${number}:${number}${number}`;
type ISODateTimeString = `${ISODateString}T${ISOTimeString}`;
export interface TimeSeriesDatum {
    date: ISODateTimeString;
    value: number;
}
export interface SimpleDatum {
    value: number;
    label: string;
}
export interface UsageDashboardData {
    auth_methods: Record<string, number>;
    kvv1_secrets: number;
    kvv2_secrets: number;
    lease_count_quotas: {
        global_lease_count_quota: {
            capacity: number;
            count: number;
            name: string;
        };
        total_lease_count_quotas: number;
    };
    namespaces: number;
    secrets_sync: number;
    pki: {
        total_issuers: number;
        total_roles: number;
    };
    replication_status: {
        dr_primary: boolean;
        dr_state: string;
        pr_primary: boolean;
        pr_state: string;
    };
    secret_engines: Record<string, number>;
}
export interface IUsageDashboardService {
    getUsageData(): Promise<UsageDashboardData>;
}
export {};
//# sourceMappingURL=index.d.ts.map