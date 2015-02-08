// Package pouch provides a set of interfaces that create a common
// interface to all backing storage media, whether it is SQL based,
// a file system, a remote storage system (i.e. S3, box), a NoSQL
// system, memcached, etc.
//
// The two core concepts are Pouches and Queries (which may soon
// be renamed to Filters, as it more appropriately describes their
// purpose).
//
// A Pouch is anything we can store data in or retrieve
// data from. In order to use a pouch for a given system, we just
// need an implementation for it (some are officialy supported,
// see the impl/ subpackage - if a given one that is needed isn't
// currently provided there, please feel free to make a pull request,
// similarly, if you can convince me that a given implementatio
// should function differently than it currently does, feel free to
// send a pull request).
//
// A Query is a filter, or predefined view, in front of a given
// Pouch. This can be used to:
//   - retrieve a certain entity without previously interactive with it
//   - implementing server side pagination, agnostic of the backing storage
//   - read only views, write only views, delete only views (permissions)
//   - add criteria that are difficult to express (cleanly) in Go
//
// NOTE
//
// This package is currently undergoing heavy development - things may
// move, things may break. If you are looking for the command line tool
// pouch (which can automagically generate go types from your backing
// storage media and can write all the necessary code to use any type
// with pouch) please see the pouch subdirectory. You can install it with:
//
//     go get github.com/ttacon/pouch/pouch
//
// Also, if you are looking for a ListAll() type method, checkout
// Query.FindEntities() - it's what you're looking for. Soon, there will be
// a nicer way to get a blank Query from a Pouch, but until then, you can
// retrieve a blank query as:
//
//     var p pouch.Pouch
//     // ... set up our pouch
//     filter := p.Offset(0)
package pouch
