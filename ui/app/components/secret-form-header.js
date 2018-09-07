import { alias } from '@ember/object/computed';
import Component from '@ember/component';
import hbs from 'htmlbars-inline-precompile';

export default Component.extend({
  key: null,
  mode: null,
  path: null,
  actionClass: null,

  title: alias('key.keyWithoutParent'),

  layout: hbs`
    <div class="consul-show-header connected">
      {{#secret-link
        mode="list"
        secret=key.parentKey
        class="back-button"
      }}
        {{i-con glyph="chevron-left" size=11}}
        Secrets
      {{/secret-link}}

      <div class="actions {{actionClass}}">
        {{yield}}
      </div>

      <div class="item-name">
        {{#if (eq mode "create") }}
          Create a secret at
          <code>
            {{#if showPrefix}}
              {{! need this to prevent a shift in the layout before we transition when saving }}
              {{#if key.isCreating}}
                {{key.initialParentKey}}
              {{else}}
                {{key.parentKey}}
              {{/if}}
            {{/if}}
          </code>
        {{/if}}

        {{#if (eq mode "edit") }}
          Edit
        {{/if}}

        <code>{{title}}</code>
      </div>
    </div>`,
});
