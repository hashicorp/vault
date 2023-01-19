// Type that comes back from expandAttributeMeta
export interface FormField {
  name: string;
  type: string;
  options: unknown;
}

export interface FormFieldGroups {
  [key: string]: Array<FormField>;
}
