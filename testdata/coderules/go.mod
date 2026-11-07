// The control of the test of the rules of code, rules_test.go at the root of learner-model: a
// module that breaks every rule, each line marked with the rules it breaks. Nothing builds it, and
// its go.sum holds placeholders for the module it requires, which does not exist.
module example.com/coderules

go 1.27.3

require example.com/dependency v1.0.0 // want require
