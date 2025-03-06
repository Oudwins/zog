---
sidebar_position: 9
---

# Performance

Zog is one of the fastest validation libraries for Go as per the [govalidbench](https://github.com/Oudwins/govalidbench) benchmarks. But to get the most performance out of zog there is some things you can do as a user:

## When possible use global schemas

One thing that makes Zog very performant is that once you have built a schema you can reuse it as many times as you want. So if you are able it is recommended not to build a new schema for each request. The difference in speed and allocations can be very drastic specially for large schemas. For reference this is the same benchmark of a large struct (as of v0.17.2) reusing the schema vs creating a new one for each request:

```bash
# Reusing the schema
BenchmarkStructComplexSuccess/Success-12                          330733              3626 ns/op             451 B/op         28 allocs/op
BenchmarkStructComplexSuccessParallel/Success-12                 1000000              1294 ns/op             476 B/op         28 allocs/op
BenchmarkStructComplexFailure/Error-12                            124696              8428 ns/op            1395 B/op         59 allocs/op
BenchmarkStructComplexFailureParallel/Error-12                    465020              3047 ns/op            1410 B/op         57 allocs/op
# Creating a new schema for each execution
BenchmarkStructComplexCreateSuccess/Success-12                    119824              8757 ns/op            9702 B/op        114 allocs/op
BenchmarkStructComplexCreateSuccessParallel/Success-12            294409              4607 ns/op            9905 B/op        114 allocs/op
BenchmarkStructComplexCreateFailure/Error-12                       80371             14069 ns/op           10652 B/op        145 allocs/op
BenchmarkStructComplexCreateFailureParallel/Error-12              192896              6993 ns/op           10793 B/op        144 allocs/op
```

If you need to build the schema on the fly and need the best performance possible I recommend you look into using `sync.Pool` to reuse the schemas.

## Use `schema.Validate` instead of `schema.Parse` when possible

For the moment parsing is slower because it needs to unmarshal the data into a map then parse it into the struct. I have quite a few ideas on how to improve `Parse` and hopefully make it as efficient as `Validate` but it will take some time. So unless you need the features that `Parse` provides I recommend you use `Validate`.

For more information on the differences read [parsing vs validation](/core-concepts/parsing-vs-validation).

## Collect issues for reuse

One of the most expensive operations in Zog is the generation of issues. We store a lot of information for the issues which is great for debugging and for generating rich custom err0r messages but can cause many allocations and be slow. A great way to mitigate this is to let Zog know that you are done using an issue. This way Zog will reuse those structs which will put less preassure on the GC. To do this Zog provides a few utility functions under the z.Issues name space:

```go
// Collects a ZogIssueMap to be reused by Zog. This will "free" the issues in the map. This can help make Zog more performant by reusing issue structs.
func (i *issueHelpers) CollectMap(issues ZogIssueMap)

// Collects a ZogIssueList to be reused by Zog. This will "free" the issues in the list. This can help make Zog more performant by reusing issue structs.
func (i *issueHelpers) CollectList(issues ZogIssueList)

// Collects a ZogIssue to be reused by Zog. This will "free" the issue. This can help make Zog more performant by reusing issue structs.
func (i *issueHelpers) Collect(issue *ZogIssue)
```
