export type HTMLElementEvent<T extends HTMLElement> = Event & {
  target: T;
  currentTarget: T;
};
