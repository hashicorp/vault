/* eslint-disable import/extensions */
import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import { withKnobs, select } from '@storybook/addon-knobs';
import notes from './toolbar-filters.md';

storiesOf('Toolbar/', module)
  .addParameters({ options: { showPanel: true } })
  .addDecorator(withKnobs())
  .add(`ToolbarFilters`, () => ({
    template: hbs`
        <h5 class="title is-5">ToolbarFilters</h5>
        <Toolbar>
          <ToolbarFilters>
            <div class="control has-icons-left">
              <input class="filter input" placeholder="Filter keys" type="text">
              <Icon @glyph="search" @size="l" class="search-icon has-text-grey-light" />
            </div>
            <div class="toolbar-separator"/>
            <div class="control">
              <input
                id="json"
                type="checkbox"
                name="json"
                class="switch is-rounded is-success is-small"
              />
              <label for="json" class="has-text-grey">JSON</label>
            </div>
          </ToolbarFilters>
        </Toolbar>
    `,
    context: {},
  }),
  {notes}
);
