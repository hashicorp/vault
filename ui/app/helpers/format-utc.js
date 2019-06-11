import { helper } from '@ember/component/helper';
import d3 from 'd3-time-format';

export function formatUtc([date, specifier]) {
  // given a UTC date, format and display it while preserving the UTC timezone.
  const parsedDate = d3.isoParse(date);
  const format = d3.utcFormat(specifier);
  return format(parsedDate);
}

export default helper(formatUtc);
