/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';

/**
 * @module ListTable::Cell
 *
 * Renders the content of a single table cell based on `column.valueType`.
 * Extracted from ListTable to eliminate the duplicated valueType branch between
 * the expandable (<B.Th>) and standard (<B.Td>) cell paths.
 *
 * Rendering is determined entirely by `column.valueType`:
 *
 * | valueType        | Config fields used                          | Renders                                      |
 * |------------------|---------------------------------------------|----------------------------------------------|
 * | "link"           | column.route, column.routeKey,              | Hds::Link::Standalone (when column.icon or   |
 * |                  | column.fallbackRouteKey,                    | row[column.iconKey] is set) or               |
 * |                  | column.routeModelKey (defaults to "id"),    | Hds::Link::Inline otherwise.                 |
 * |                  | column.iconKey, column.icon                 | Resolves route: routeKey → fallbackRouteKey  |
 * |                  |                                             | → route (static).                            |
 * | "icon-text"      | column.iconKey, column.icon (fallback       | Hds::Icon + optional Hds::TooltipButton      |
 * |                  | "lock"), column.tooltipKey, column.tooltip, | wrapping the icon + plain text span.         |
 * |                  | column.textKey (falls back to column.key)   | Tooltip shown when tooltipKey or tooltip set.|
 * | "copy"           | column.key (value)                          | Hds::Copy::Button (icon-only) + plain text.  |
 * |                  |                                             | Renders nothing when value is falsy.         |
 * | "status"         | column.statusMap (value → badge color),     | Hds::Badge with @text, @color, and @icon.    |
 * |                  | column.statusIconMap (value → badge icon)   | Unrecognized values render a neutral badge.  |
 * | "date-time"      | column.key (ISO date string value)          | <time> wrapping {{date-format}} output;      |
 * |                  |                                             | format: "MMM dd, yyyy h:mm a".               |
 * |                  |                                             | Renders nothing when value is falsy.         |
 * | "custom" /       | —                                           | Yields <:customTableItem as |row col val|>   |
 * | customTableItem  |                                             | back to the consuming template.              |
 * | (anything else)  | —                                           | Plain text with word wrapping.               |
 *
 * @param {object} column  - Column definition from ListViewConfig.columns
 * @param {object} row     - The full row data object
 * @param {unknown} value  - Pre-resolved cell value: get(row, column.key)
 */

interface Args {
  column: Record<string, unknown>;
  row: Record<string, unknown>;
  value: unknown;
}

export default class ListTableCell extends Component<Args> {
  isObject = (val: unknown) => typeof val === 'object' && val !== null;
}
