// Package learnermodel is NovaTrail's learner model: it estimates how well a learner has mastered
// each skill of a module from the evidence of finished attempts, and computes each skill's state.
//
// It is a library with no input or output of its own. The api service hands it the evidence and
// the module's skill graph, keeps the evidence log and the state the library computes from it,
// and passes the time in as a number; the engine reads that state to pick what comes next. The
// library picks no lesson or task, gives no verdict, calls no AI model and reads no module
// package or database.
//
// The state is a fold of the evidence log, so appending a record gives, bit for bit, the state
// recomputed from the whole log, on any build of one version of this module, WebAssembly
// included. The code imports nothing but Go's standard library, and only the packages its rules
// allow.
//
// What the library computes and why is in docs/architecture.md, section 3, and in
// docs/decisions.md of this repository.
package learnermodel
