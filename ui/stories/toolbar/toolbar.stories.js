/* eslint-disable import/extensions */
import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import { withKnobs, select } from '@storybook/addon-knobs';
import notes from './toolbar.md';

storiesOf('Toolbar/', module)
  .addParameters({ options: { showPanel: true } })
  .addDecorator(withKnobs())
  .add(`Toolbar`, () => ({
    template: hbs`
        <h5 class="title is-5">Toolbar</h5>
        <Toolbar>
          {{#if (eq example "List")}}
            <ToolbarFilters>
              <div class="control has-icons-left">
                <input class="filter input" placeholder="Filter keys" type="text">
                <Icon @glyph="search" class="search-icon has-text-grey-light hs-icon-l" />
              </div>
            </ToolbarFilters>
            <ToolbarActions>
              <ToolbarLink
                @type="add"
                @params={{array "#"}}
              >
                Add item
              </ToolbarLink>
            </ToolbarActions>
          {{/if}}
          {{#if (eq example "Show")}}
            <ToolbarActions>
              <ToolbarLink
                @params={{array "#"}}
              >
                Delete
              </ToolbarLink>
              <ToolbarLink
                @params={{array "#"}}
              >
                Edit
              </ToolbarLink>
            </ToolbarActions>
          {{/if}}
          {{#if (eq example "Code editor")}}
            <ToolbarFilters>
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
            <ToolbarActions>
              <ToolbarLink
                @params={{array "#"}}
              >
                Copy
              </ToolbarLink>
            </ToolbarActions>
          {{/if}}
        </Toolbar>
    `,
    context: {
      example: select('Example', ['List', 'Show', 'Code editor'], 'List'),
    },
  }),
  {notes}
);
