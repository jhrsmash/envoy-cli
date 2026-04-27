// Package envfile provides the Archive feature for envoy-cli.
//
// Archive allows users to save named, timestamped snapshots of an env map to
// disk. Each archive is stored as a JSON file under a configurable directory,
// making it easy to restore or diff historical environment states.
//
// # Core functions
//
//   - Archive(env, label, dir)       – save a labelled snapshot to dir
//   - LoadArchive(label, dir)        – load a previously saved snapshot
//   - ListArchives(dir)              – enumerate all archive labels in dir
//
// # Formatting
//
//   - FormatArchiveResult(r)         – summary of an archive write
//   - FormatArchiveEntry(e)          – human-readable view of a loaded entry
//   - FormatArchiveList(labels)      – list of available archive labels
package envfile
