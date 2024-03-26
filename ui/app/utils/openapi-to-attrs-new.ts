// /**
//  * Copyright (c) HashiCorp, Inc.
//  * SPDX-License-Identifier: BUSL-1.1
//  */

// import { attr } from '@ember-data/model';
// import { camelize, capitalize } from '@ember/string';

// interface OpenApiProp {
//   description: string;
//   items: { type: string };
//   type: string;
//   format: string;
//   isId: boolean;
//   deprecated: boolean;
//   enum: string[];
//   'x-vault-displayAttrs': {
//     name: string;
//     value: string;
//     group: string;
//     sensitive: boolean;
//     editType: string;
//     description: string;
//   };
// }
// interface Props {
//   [key: string]: OpenApiProp;
// }
// interface AttrDefn {
//   editType: string;
//   helpText: string;
//   possibleValues: string[];
//   fieldValue: string | null;
//   fieldGroup: string;
//   readOnly: boolean;
//   defaultValue: unknown;
//   sensitive?: boolean;
//   label?: string;
//   type: string;
// }
// export const expandOpenApiProps = function (props: Props) {
//   const attrs = {};
//   // expand all attributes
//   for (const propName in props) {
//     const prop = props[propName] as OpenApiProp;
//     let { description, items, type, format, isId, deprecated } = prop;
//     if (deprecated === true) {
//       continue;
//     }
//     let {
//       name,
//       value,
//       group,
//       sensitive,
//       editType,
//       description: displayDescription,
//     } = prop['x-vault-displayAttrs'] || {};

//     if (type === 'integer') {
//       type = 'number';
//     }

//     if (displayDescription) {
//       description = displayDescription;
//     }

//     editType = editType || type;

//     if (format === 'seconds' || format === 'duration') {
//       editType = 'ttl';
//     } else if (items) {
//       editType = items.type + capitalize(type);
//     }

//     let attrType = 'string';
//     // ttls write as a string and read as a number
//     // so setting type on them runs the wrong transform
//     if (editType !== 'ttl' && type !== 'array') {
//       attrType = type;
//     }

//     const attrDefn: AttrDefn = {
//       type: attrType,
//       editType,
//       helpText: description,
//       possibleValues: prop.enum,
//       fieldValue: isId ? 'mutableId' : null,
//       fieldGroup: group || 'default',
//       readOnly: isId,
//       defaultValue: value || null,
//     };

//     if (type === 'object' && !!value) {
//       attrDefn.defaultValue = () => {
//         return value;
//       };
//     }

//     if (sensitive) {
//       attrDefn.sensitive = true;
//     }

//     // only set a label if we have one from OpenAPI
//     // otherwise the propName will be humanized by the form-field component
//     if (name) {
//       attrDefn.label = name;
//     }

//     // loop to remove empty vals
//     for (const attrProp in attrDefn) {
//       if (attrDefn[attrProp as keyof typeof attrDefn] == null) {
//         delete attrDefn[attrProp as keyof typeof attrDefn];
//       }
//     }
//     const attrKey = camelize(propName as string);
//     attrs[ as keyof typeof attrDefn] = attrDefn;
//   }
//   return attrs;
// };

// export const combineAttributes = function (oldAttrs, newProps) {
//   const newAttrs = {};
//   const newFields = [];
//   console.log('CHECK', oldAttrs, newProps);
//   if (oldAttrs) {
//     // make sure it returns something
//     oldAttrs.forEach(function (value, name) {
//       if (newProps[name]) {
//         newAttrs[name] = attr(newProps[name].type, { ...newProps[name], ...value.options });
//       } else {
//         newAttrs[name] = attr(value.type, value.options);
//       }
//     });
//   }
//   for (const prop in newProps) {
//     if (newAttrs[prop]) {
//       continue;
//     } else {
//       console.log('newProps[name]', prop, newProps[prop]);
//       newAttrs[prop] = attr(newProps[prop].type || 'string', newProps[prop]);
//       newFields.push(prop);
//     }
//   }
//   return { attrs: newAttrs, newFields };
// };

// export const combineFieldGroups = function (currentGroups, newFields, excludedFields) {
//   let allFields = [];
//   for (const group of currentGroups) {
//     const fieldName = Object.keys(group)[0];
//     allFields = allFields.concat(group[fieldName]);
//   }
//   const otherFields = newFields.filter((field) => {
//     return !allFields.includes(field) && !excludedFields.includes(field);
//   });
//   if (otherFields.length) {
//     currentGroups[0].default = currentGroups[0].default.concat(otherFields);
//   }

//   return currentGroups;
// };
