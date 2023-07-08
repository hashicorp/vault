/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

export type HTMLElementEvent<T extends HTMLElement> = Event & {
  target: T;
  currentTarget: T;
};
