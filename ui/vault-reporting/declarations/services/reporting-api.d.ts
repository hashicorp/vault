/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
import Service from '@ember/service';
import { Configuration, VaultReportingServiceApi } from '@hashicorp/cloud-vault-reporting-typescript';
export default class ReportingApiService extends Service {
    /**
     * Attempt to look up the API service from host application if it exists. Using service decorator would error so this allows falling back if it doesn't.
     */
    get api(): {
        config: (basePath: string) => Configuration;
    } | undefined;
    /**
     * Attempt to look up the user context from host application if it exists. Using service decorator would error so this allows falling back if it doesn't.
     */
    get userContext(): {
        organizationId?: string;
        projectId?: string;
    } | undefined;
    /**
     * Attempt to get the organizationId and projectId from user context in host application.
     * Fall back to default values if not available.
     */
    get organizationId(): string;
    get projectId(): string;
    /**
     * Attempt to get the API configuration from host application. Fall back to default otherwise.
     */
    get config(): Configuration;
    reporting: VaultReportingServiceApi;
}
//# sourceMappingURL=reporting-api.d.ts.map