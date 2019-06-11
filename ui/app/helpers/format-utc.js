import { helper } from '@ember/component/helper';
import d3 from 'd3-time-format';

export function formatUtc([date, specifier]) {
  const parsedDate = d3.isoParse(date);
  const formatter = d3.utcFormat(specifier);
  return formatter(parsedDate);
}

export default helper(formatUtc);
