
# FormFieldGroupsLoop
FormFieldGroupsLoop components loop through the groups set on a model and display them either as default or behind toggle components.

| Param | Type | Default | Description |
| --- | --- | --- | --- |
| model | <code>class</code> |  | The routes model class. |
| mode | <code>string</code> |  | "create" or "update" used to hide the name form field. TODO: not ideal, would prefer to disable it to follow new design patterns. |
| [modelValidations] | <code>function</code> |  | Passed through to formField. |
| [showHelpText] | <code>boolean</code> |  | Passed through to formField. |
| [groupName] | <code>string</code> | <code>&quot;\&quot;fieldGroups\&quot;&quot;</code> | option to override key on the model where groups are located |

**Example**  
```hbs preview-template
<FormFieldGroupsLoop @model={{this.model}} @mode={{if @model.isNew "create" "update"}}/>
```
