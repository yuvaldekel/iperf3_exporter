# Changelog

## [2.0.0](https://github.com/yuvaldekel/iperf3_exporter/compare/v1.3.1...v2.0.0) (2026-02-20)


### ⚠ BREAKING CHANGES

* Several user-facing changes require updates:
    - Docker image moved from edgard/iperf3-exporter to ghcr.io/edgard/iperf3_exporter
    - Docker container now runs as non-root user
    - Removed go-kit/log dependency and changed logging implementation
    - Modified HTTP server implementation with prometheus/exporter-toolkit

### Features

* add bind parameter support and fix bandwidth calculation docs ([#38](https://github.com/yuvaldekel/iperf3_exporter/issues/38)) ([7457e9e](https://github.com/yuvaldekel/iperf3_exporter/commit/7457e9ecab88cc68584960842afbab0cdaf1c2e8)), closes [#34](https://github.com/yuvaldekel/iperf3_exporter/issues/34) [#35](https://github.com/yuvaldekel/iperf3_exporter/issues/35)
* add pprof handlers for performance profiling ([c9e4180](https://github.com/yuvaldekel/iperf3_exporter/commit/c9e41801291063c7cf40e63564101b44a7870b78))
* add support for UDP mode ([e7fd775](https://github.com/yuvaldekel/iperf3_exporter/commit/e7fd775507ba15f2cd408ec818b86db24408749c))
* restructure project to standard Go layout and modernize ([18fa4cc](https://github.com/yuvaldekel/iperf3_exporter/commit/18fa4ccea20ed47d42ca35fc149c64f2ebc6fea4))


### Bug Fixes

* add configuration validation and move landing page to it's own file ([34438ee](https://github.com/yuvaldekel/iperf3_exporter/commit/34438ee072c16b09546a6b512276d8d3ca1ec102))
* address multiple potential bugs ([78510be](https://github.com/yuvaldekel/iperf3_exporter/commit/78510be206fa4880c943d08556711fc8c4342825))
* correct release-please config filename and add manifest ([#41](https://github.com/yuvaldekel/iperf3_exporter/issues/41)) ([eeebf87](https://github.com/yuvaldekel/iperf3_exporter/commit/eeebf87aaa71f36be22ee096f6091fffca8fc8b8))
* CVE-2022-21698 and dependencies update ([40fb832](https://github.com/yuvaldekel/iperf3_exporter/commit/40fb8320acf6e58f97ab858c53e02782c59adf30))
* improve code organization and error handling ([d9150fd](https://github.com/yuvaldekel/iperf3_exporter/commit/d9150fd1d964a789a030bad396731f7182ff8402))
* metrics collection and timeout issues ([7428bc8](https://github.com/yuvaldekel/iperf3_exporter/commit/7428bc880d8e9e4b389f05875eb468175fccd852))
* race condition running iperf3 command and associated tests ([c39cc93](https://github.com/yuvaldekel/iperf3_exporter/commit/c39cc93404b34d987f8a2db2f21560e8d17b90ba))
* remove GitHub repository link from index handler and update test expectations ([fdd11be](https://github.com/yuvaldekel/iperf3_exporter/commit/fdd11be2a33cf2a853d5c6ba02759d50aa674b30))
* respect --iperf3.timeout as upper limit when Prometheus header present ([#39](https://github.com/yuvaldekel/iperf3_exporter/issues/39)) ([352f215](https://github.com/yuvaldekel/iperf3_exporter/commit/352f21540be81ff493a2be14d1b5a5d95f018fd4)), closes [#37](https://github.com/yuvaldekel/iperf3_exporter/issues/37)
* simplify metrics, remove old tests ([7164e64](https://github.com/yuvaldekel/iperf3_exporter/commit/7164e640b5b97ffc4d3a8c4d16aefb018c67a4d3))
* version display missing ([478c489](https://github.com/yuvaldekel/iperf3_exporter/commit/478c48933d4fa5711fcbd1e4a53f5831745ed5d1))
