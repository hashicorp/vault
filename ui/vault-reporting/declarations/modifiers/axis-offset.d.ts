/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
/**
 * By default the axis elements are outside of the bounds of the svg and rely on the containing element having
 * enough padding to compensate. A fixed padding is not flexible to varied width of axis labels.
 *
 * This modifier is used to pad compensate for the width of the axis element. It also returns a value
 * that can be used to set the chart width based on the offset amount.
 */
declare const _default: import("ember-modifier").FunctionBasedModifier<{
    Args: {
        Positional: [(offset: number) => unknown, additionalPadding?: number | undefined];
        Named: import("ember-modifier/-private/signature").EmptyObject;
    };
    Element: SVGElement;
}>;
export default _default;
//# sourceMappingURL=axis-offset.d.ts.map