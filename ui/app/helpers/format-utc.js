import { helper } from '@ember/component/helper';
import d3 from 'd3-time-format';

export function formatUtc([date, specifier]) {
  // given a date, format and display it as UTC.
  const format = d3.utcFormat(specifier);
  const parse = d3.utcParse('%Y-%m-%dT%H:%M:%SZ');

  // if a date isn't already in UTC, fallback to isoParse to convert it to UTC
  const parsedDate = parse(date) || d3.isoParse(date);

  return format(parsedDate);
}

export default helper(formatUtc);
