# HyperLogLog - an algorithm for approximating the number of distinct elements

[![GoDoc](https://godoc.org/github.com/axiomhq/hyperloglog?status.svg)](https://godoc.org/github.com/axiomhq/hyperloglog) [![Go Report Card](https://goreportcard.com/badge/github.com/axiomhq/hyperloglog)](https://goreportcard.com/report/github.com/axiomhq/hyperloglog) [![CircleCI](https://circleci.com/gh/axiomhq/hyperloglog/tree/master.svg?style=svg)](https://circleci.com/gh/axiomhq/hyperloglog/tree/master)

An improved version of [HyperLogLog](https://en.wikipedia.org/wiki/HyperLogLog) for the count-distinct problem, approximating the number of distinct elements in a multiset. This implementation offers enhanced performance, flexibility, and simplicity while maintaining accuracy.

## Note on Implementation History

The initial version of this work (tagged as v0.1.0) was based on ["Better with fewer bits: Improving the performance of cardinality estimation of large data streams - Qingjun Xiao, You Zhou, Shigang Chen"](http://cse.seu.edu.cn/PersonalPage/csqjxiao/csqjxiao_files/papers/INFOCOM17.pdf). However, the current implementation has evolved significantly from this original basis, notably moving away from the tailcut method.

## Current Implementation

The current implementation is based on the LogLog-Beta algorithm, as described in:

["LogLog-Beta and More: A New Algorithm for Cardinality Estimation Based on LogLog Counting"](https://arxiv.org/pdf/1612.02284) by Jason Qin, Denys Kim, and Yumei Tung (2016).

Key features of the current implementation:
* **Metro hash** used instead of xxhash
* **Sparse representation** for lower cardinalities (like HyperLogLog++)
* **LogLog-Beta** for dynamic bias correction across all cardinalities
* **8-bit registers** for convenience and simplified implementation
* **Order-independent insertions and merging** for consistent results regardless of data input order
* **Removal of tailcut method** for a more straightforward approach
* **Flexible precision** allowing for 2^4 to 2^18 registers

This implementation is now more straightforward, efficient, and flexible, while remaining backwards compatible with previous versions. It provides a balance between precision, memory usage, speed, and ease of use.

## Precision and Memory Usage

This implementation allows for creating HyperLogLog sketches with arbitrary precision between 2^4 and 2^18 registers. The memory usage scales with the number of registers:

* Minimum (2^4 registers): 16 bytes
* Default (2^14 registers): 16 KB
* Maximum (2^18 registers): 256 KB

Users can choose the precision that best fits their use case, balancing memory usage against estimation accuracy.

## Note
A big thank you to Prof. Shigang Chen and his team at the University of Florida who are actively conducting research around "Big Network Data".

## Contributing

Kindly check our [contributing guide](https://github.com/axiomhq/hyperloglog/blob/main/Contributing.md) on how to propose bugfixes and improvements, and submitting pull requests to the project

## License

&copy; Axiom, Inc., 2024

Distributed under MIT License (`The MIT License`).

See [LICENSE](LICENSE) for more information.