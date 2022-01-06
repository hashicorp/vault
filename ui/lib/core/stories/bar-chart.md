<!--THIS FILE IS AUTO GENERATED. This file is generated from JSDoc comments in lib/core/addon/components/bar-chart.js. To make changes, first edit that file and run "yarn gen-story-md bar-chart" to re-generate the content.-->

## BarChart
BarChart components are used to display data in the form of a stacked bar chart, with accompanying legend and tooltip. Anything passed into the block will display in the top right of the chart header.

**Params**

| Param | Type | Default | Description |
| --- | --- | --- | --- |
| title | <code>string</code> |  | title of the chart |
| mapLegend | <code>array</code> |  | array of objects with key names 'key' and 'label' for the map legend |
| dataset | <code>object</code> |  | dataset for the chart |
| [description] | <code>string</code> |  | description of the chart |
| [labelKey] | <code>string</code> | <code>&quot;label&quot;</code> | labelKey is the key name in the dataset passed in that corresponds to the value labeling the y-axis |
| [onClick] | <code>function</code> |  | takes function from parent and passes it to click event on data bars |

**Example**
  
```js
<BarChartComponent @title="Top 10 Namespaces" @description="Each namespace's client count includes clients in child namespaces." @labelKey="namespace_path" @dataset={{this.testData}} @mapLegend={{ array (hash key="non_entity_tokens" label="Active direct tokens") (hash key="distinct_entities" label="Unique Entities") }} @onClick={{this.onClick}} >
   <button type="button" class="link">
     Export all namespace data
   </button>/>
</BarChartComponent>

 mapLegendSample = [{
    key: "api_key_for_label",
    label: "Label Displayed on Legend"
  }]
```

**See**

- [Uses of BarChart](https://github.com/hashicorp/vault/search?l=Handlebars&q=BarChart+OR+bar-chart)
- [BarChart Source Code](https://github.com/hashicorp/vault/blob/master/ui/lib/core/addon/components/bar-chart.js)

---
