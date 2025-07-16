/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/**
 * @module FilterInputExplicit
 *
 * @description FilterInputExplicit component is a child component to show filter input.
 * It also handles the filtering actions of roles.
 *
 * @example
 * <FilterInputExplicit
 *   @query={{this.pageFilter}}
 *   @placeholder="Search"
 *   @handleSearch={{this.handleSearch}}
 *   @handleInput={{this.handleInput}}
 *   @handleKeyDown={{this.handleKeyDown}}
 * />
 *
 * @param {string} query - value of queryParam, such as pageFilter
 * @param {string} placeholder - placeholder for the input field
 * @param {function} handleSearch - callback function to handle search
 * @param {function} handleInput - callback function to handle input
 * @param {function} handleKeyDown - callback function to handle keydown
 */
export { default } from 'core/components/filter-input-explicit';
