# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.2] - 2022-01-03
### Fixed
- Critical logic error has been fixed where non-range verse numbers would result in an out-of-bounds index access

## [1.1.1] - 2022-01-03
### Changed
- Implemented bitflags for functions using the [bitflags](https://github.com/bitflags/bitflags) Rust library, hopefully
increasing the performance, but also streamlining the code

### Fixed
- Fixed a bug where the end verse would be read with an offset of minus two (e.g. "6:20" would get you verses 6 to 18,
not 6 to 20 as expected by the user)

## [1.1.0] - 2022-01-01
### Added
- Configuration flag (`-c, --colour`) to colour output annotations
- Configuration flag (`-v, --verbose`) to output warnings about user input
- Parallel computation using Rust library Rayon, dramatically increasing performance

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
