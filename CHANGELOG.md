# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.0] - 2022-01-01
### Added
- Configuration flag (`-c, --colour`) to colour output annotations
- Configuration flag (`-v, --verbose`) to output warnings about user input
- Parallel computation using Rust library Rayon. Increases printing of the Bible from ~500 milliseconds to ~350
  milliseconds.

### Changed
- Rewritten from scratch from Go in Rust
- Due to parallel computation, output is no longer in exact canonical order

### Fixed
- Fixed an out-of-bounds index error when providing too large verse ranges

## [1.0.0] - 2021-09-25
### Added
- Printing the entire Bible
- Printing a book from the Bible
- Printing a chapter from the Bible
- Printing a specific verse, or set of verses, from the Bible
