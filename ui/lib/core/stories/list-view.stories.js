import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import { withKnobs, select } from '@storybook/addon-knobs';
import notes from './list-view.md';

import ArrayProxy from '@ember/array/proxy';

let filtered = ArrayProxy.create({ content: [] });
filtered.set('meta', {
  lastPage: 1,
  currentPage: 1,
  total: 100,
});

let paginated = ArrayProxy.create({
  content: [{ id: 'middle' }, { id: 'of' }, { id: 'the' }, { id: 'list' }],
});
paginated.set('meta', {
  lastPage: 10,
  currentPage: 4,
  total: 100,
});

let options = {
  list: [{ id: 'one' }, { id: 'two' }],
  empty: [],
  filtered,
  paginated,
};

storiesOf('ListView/', module)
  .addParameters({ options: { showPanel: true } })
  .addDecorator(withKnobs())
  .add(
    `ListView`,
    () => ({
      template: hbs`
				<h5 class="title is-5">{{title}}</h5>
				<ListView @items={{items}} @itemNoun={{or noun "role"}} @paginationRouteName="vault" as |list|>
					{{#if list.empty}}
						<list.empty @title="No roles here" />
					{{else if list.item}}
						<div class="box is-marginless">
							{{list.item.id}}
						</div>
					{{else}}
						<div class="box">There aren't any items in this filter</div>
					{{/if}}
				</ListView>
	`,
      context: {
        title: 'ListView',
        items: select('items', options, options['list']),
      },
    }),
    { notes }
  );
