/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { ResponseError } from '@hashicorp/vault-client-typescript';
import { action } from '@ember/object';
import routerLookup from 'core/utils/router-lookup';

import type Owner from '@ember/owner';
import type ApiService from 'vault/services/api';
import type RouterService from '@ember/routing/router-service';

interface Args {
  error: ResponseError | PageError;
}

type ApplicationError = {
  errorURL?: string; // Added by application route error() method
  httpStatus: number; // Added by application adapter handleRequest() method
  path: string; // Added by application adapter handleRequest() method
};

type NotFoundRouteError = {
  httpStatus: number;
  path: string;
};

type ParsedApiServiceError = {
  errorURL?: string; // Added by the application route error() method
  message: string; //The API service sets message from response errors when present
  path?: string;
  response?: { errors: string[] };
  status?: number;
};

// Can be removed when Ember Data is no more
type EmberDataAdapterError = {
  errors: string[];
  errorURL?: string; // Added by the application route error() method
  httpStatus: number;
  message: string;
  path: string;
};

type PageError = ApplicationError | NotFoundRouteError | ParsedApiServiceError | EmberDataAdapterError;

enum DefaultErrorMessages {
  NOT_FOUND = 'Sorry, we were unable to find any content at',
  NOT_AUTHORIZED = 'You are not authorized to access content at',
  UNKNOWN = 'A problem has occurred. Check the Vault logs or console for more details.',
}

export default class PageErrorComponent extends Component<Args> {
  @service declare readonly api: ApiService;

  // Since our error handling is so inconsistent, this component is set up to handle an undefined error
  @tracked declare error: PageError | undefined;

  constructor(owner: Owner, args: Args) {
    super(owner, args);
    this.unpackError();
  }

  get displayUrl() {
    // "Path" is the API path for which the request failed, show that if set.
    // Otherwise errorURL is the URL when the request failed.
    // Finally, fallback to the current URL.
    const url = this.error && 'errorURL' in this.error ? this.error.errorURL : this.router.currentURL;
    return this.error?.path || url;
  }

  // See typescript definitions above. Errors have a slightly different shape depending on
  // where/how they are thrown.
  get errorCode() {
    if (!this.error) {
      return '';
    }
    if ('httpStatus' in this.error) {
      return this.error.httpStatus;
    }
    if ('status' in this.error) {
      return this.error.status;
    }
    return '';
  }

  get shouldRenderUrl() {
    const { NOT_FOUND, NOT_AUTHORIZED } = DefaultErrorMessages;
    return this.message === NOT_FOUND || this.message === NOT_AUTHORIZED;
  }

  // In most cases, we want to render the message returned by the API, when available.
  // However, if it is a generic "permission denied" message or from an Ember Data adapter
  // then the default messages provide more information and context.
  // When Ember Data is removed, that part of the conditional can be removed.
  get message(): string | DefaultErrorMessages {
    const message = this.error && 'message' in this.error ? this.error.message : undefined;
    if (!message || message?.includes('permission denied') || message?.includes('Ember Data')) {
      return this.statusDefaults.message;
    }
    return message;
  }

  get router() {
    return routerLookup(this) as RouterService;
  }

  get statusDefaults() {
    switch (this.errorCode) {
      case 404:
        return {
          icon: 'alert-circle',
          title: 'Not found',
          message: DefaultErrorMessages.NOT_FOUND,
        };
      case 403:
        return {
          icon: 'skip',
          title: 'Not authorized',
          message: DefaultErrorMessages.NOT_AUTHORIZED,
        };
      default:
        return {
          icon: 'alert-circle',
          title: 'Error',
          message: DefaultErrorMessages.UNKNOWN,
        };
    }
  }

  @action
  async unpackError() {
    if (this.args.error instanceof ResponseError) {
      // Pass an empty string for the fallback message because this component has its own fallback message handling
      const { status, path, message } = await this.api.parseError(this.args.error, '');
      this.error = {
        status,
        path,
        message,
      };
    } else {
      this.error = this.args.error;
    }
  }
}
