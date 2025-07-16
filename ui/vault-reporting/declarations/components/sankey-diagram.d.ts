/**
 * Copyright (c) HashiCorp, Inc.
 */

import Component from '@glimmer/component';
interface SankeyDiagramArgs {
    data: {
        nodes: {
            name: string;
        }[];
        links: {
            source: number;
            target: number;
            value: number;
        }[];
    };
}
export default class SankeyDiagramComponent extends Component<SankeyDiagramArgs> {
    renderSankey(element: HTMLElement): void;
}
export {};
//# sourceMappingURL=sankey-diagram.d.ts.map