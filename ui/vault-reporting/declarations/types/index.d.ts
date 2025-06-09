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
export declare enum REPLICATION_ENABLED_STATE {
    PRIMARY = "primary",
    SECONDARY = "secondary",
    BOOTSTRAPPING = "bootstrapping"
}
export declare const REPLICATION_DISABLED_STATE = "disabled";
export interface UsageDashboardData {
    authMethods: Record<string, number>;
    leasesByAuthMethod: Record<string, number>;
    kvv1Secrets: number;
    kvv2Secrets: number;
    leaseCountQuotas: {
        globalLeaseCountQuota: {
            capacity: number;
            count: number;
            name: string;
        };
        totalLeaseCountQuotas: number;
    };
    namespaces: number;
    secretSync: {
        totalDestinations: number;
    };
    pki: {
        totalIssuers: number;
        totalRoles: number;
    };
    replicationStatus: {
        drPrimary: boolean;
        drState: REPLICATION_ENABLED_STATE | typeof REPLICATION_DISABLED_STATE;
        prPrimary: boolean;
        prState: REPLICATION_ENABLED_STATE | typeof REPLICATION_DISABLED_STATE;
    };
    secretEngines: Record<string, number>;
}
export type getUsageDataFunction = () => Promise<UsageDashboardData>;
export {};
//# sourceMappingURL=index.d.ts.map