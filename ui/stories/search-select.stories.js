/* eslint-disable import/extensions */
import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import notes from './search-select.md';
import { withKnobs, text, select } from '@storybook/addon-knobs';

const onChange = (value) => alert(`New value is "${value}"`);
const models = ["identity/groups"];

storiesOf('SearchSelect/', module)
  .addParameters({ options: { showPanel: true } })
  .addDecorator(withKnobs({ escapeHTML: false }))
  .add(`SearchSelect`, () => ({
    template: hbs`
      <h5 class="title is-5">Search Select</h5>
      <SearchSelect 
        @id="groups" 
        @models={{models}} 
        @onChange={{onChange}} 
        @inputValue={{inputValue}}
        @label={{label}}
        @fallbackComponent="string-list" 
        @staticOptions={{staticOptions}}/>
    `,
    context: {
      label: text("Label", "Group IDs"),
      helpText: text("Help Tooltip Text", "Group IDs to associate with this entity"),
      inputValue: [],
      models: models,
      onChange: onChange,
      staticOptions: [{ name: "my-group", id: "123dsafdsarf" }, { name: "my-other-group", id: "45ssadd435" }, { name: "example-1", id: "5678" }, { name: "group-2", id: "gro09283" }],
    },
  }),
    { notes }
  );
