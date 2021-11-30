import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import { object, text, withKnobs } from '@storybook/addon-knobs';
import notes from './bar-chart.md';

const dataset = [
  {
    namespace_id: 'root',
    namespace_path: 'root',
    counts: {
      distinct_entities: 268,
      non_entity_tokens: 985,
      clients: 1253,
    },
  },
  {
    namespace_id: 'O0i4m',
    namespace_path: 'top-namespace',
    counts: {
      distinct_entities: 648,
      non_entity_tokens: 220,
      clients: 868,
    },
  },
  {
    namespace_id: '1oihz',
    namespace_path: 'anotherNamespace',
    counts: {
      distinct_entities: 547,
      non_entity_tokens: 337,
      clients: 884,
    },
  },
  {
    namespace_id: '1oihz',
    namespace_path: 'someOtherNamespaceawgagawegawgawgawgaweg',
    counts: {
      distinct_entities: 807,
      non_entity_tokens: 234,
      clients: 1041,
    },
  },
];

const flattenData = () => {
  return dataset.map(d => {
    return {
      label: d['namespace_path'],
      non_entity_tokens: d['counts']['non_entity_tokens'],
      distinct_entities: d['counts']['distinct_entities'],
      total: d['counts']['clients'],
    };
  });
};

storiesOf('BarChart', module)
  .addParameters({ options: { showPanel: true } })
  .addDecorator(withKnobs())
  .add(
    `BarChart`,
    () => ({
      template: hbs`
      <h5 class="title is-5">Bar Chart</h5>
      
      <p> <code>dataset</code> is passed to a function in the parent to format it appropriately for the chart. Any data passed should be flattened (not nested).</p>
      <p> The legend typically displays within the bar chart border, below the second grey divider. There is also a tooltip that pops up when hovering over the data bars and overflowing labels. Gotta love storybook :) </p>
      <div class="chart-container" style="margin-top:24px; max-width:750px; max-height:500px;" >
        <BarChart 
          @title={{title}}
          @description={{description}}
          @dataset={{dataset}}
          @mapLegend={{array
          (hash key="non_entity_tokens" label="Active direct tokens")
          (hash key="distinct_entities" label="Unique Entities")}}
        >    
          <button type="button" class="link">
          Export all namespace data
          </button>
        </BarChart>
      <br>
      <h6 class="title is-6">Legend:</h6>
        <svg class="legend">
          <circle cx="60%" cy="10%" r="6" style="fill: rgb(191, 212, 255);"></circle>
          <text x="62%" y="10%" alignment-baseline="middle" style="font-size: 0.8rem;">Active direct tokens</text>
          <circle cx="80%" cy="10%" r="6" style="fill: rgb(138, 177, 255);"></circle>
          <text x="82%" y="10%" alignment-baseline="middle" style="font-size: 0.8rem;">Unique Entities</text></svg>
      </div>
    `,
      context: {
        title: text('title', 'Top Namespaces'),
        description: text(
          'description',
          'Each namespaces client count includes clients in child namespaces.'
        ),
        dataset: object('dataset', flattenData()),
      },
    }),
    { notes }
  );
