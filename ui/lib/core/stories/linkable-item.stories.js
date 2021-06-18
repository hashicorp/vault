import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import { withKnobs } from '@storybook/addon-knobs';

storiesOf('LinkableItem', module)
  .addParameters({ options: { showPanel: true } })
  .addDecorator(withKnobs())
  .add(`LinkableItem`, () => ({
    template: hbs`
      <h5 class="title is-5">Linkable Item</h5>
      <LinkableItem @title="My title" @accessor="v2 some-accessor" @description="Some description here" />
    `,
    context: {},
  }))
  .add(`LinkableItem with stuff`, () => ({
    template: hbs`
      <h5 class="title is-5">Linkable Item</h5>
      <LinkableItem @title="My title" @accessor="v2 some-accessor"
      @description="Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nam volutpat vulputate lacus sit amet lobortis. Nulla fermentum porta consequat. Mauris porttitor libero nibh, ac facilisis ex molestie non. Nulla dolor est, pharetra et maximus vel, varius eu augue. Maecenas eget nisl convallis, vehicula massa quis, pharetra justo. Praesent porttitor arcu at gravida dignissim. Vestibulum condimentum, risus a fermentum pulvinar, enim massa venenatis velit, a venenatis leo sem eget dolor. Morbi convallis dui sit amet egestas commodo. Nulla et ultricies leo. Integer ac semper ipsum. Donec in fermentum velit, et suscipit mauris. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur ultrices risus urna, at commodo tellus aliquet eu. Praesent ullamcorper ultrices arcu id euismod. Etiam bibendum nisi quis tortor ultrices fermentum. "
      >
        Stuff
      </LinkableItem>
    `,
    context: {},
  }));
