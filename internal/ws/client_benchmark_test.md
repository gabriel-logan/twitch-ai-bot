# Benchmark Analysis: `config.Env` by Value vs Pointer

## Objective

This benchmark was created to evaluate the performance impact of passing `config.Env` as:

* a **value** (`config.Env`)
* a **pointer** (`*config.Env`)

The goal is to determine the most efficient approach for the current project architecture, especially in hot execution paths such as the WebSocket bot runtime.

---

# Tested Structure

The benchmark uses the real `config.Env` structure containing multiple fields:

* strings
* slices
* integers
* durations

This makes the test representative of the actual production environment.

---

# Benchmarks Implemented

## 1. Environment Creation

Compares the cost of creating:

* `setupEnvValue() config.Env`
* `setupEnvPtr() *config.Env`

---

## 2. Passing Environment to Functions

Compares:

* passing the entire struct by value
* passing only the pointer reference

---

## 3. Accessing All Struct Fields

Measures access to every field inside the struct to evaluate read performance after the object is already in memory.

---

# Benchmark Results

## Run 1

```bash
BenchmarkSetupEnvValue-8           1000000000   0.3071 ns/op   0 B/op   0 allocs/op
BenchmarkSetupEnvPtr-8             1000000000   0.3111 ns/op   0 B/op   0 allocs/op
BenchmarkPassEnvValue-8              78891510  14.58 ns/op    0 B/op   0 allocs/op
BenchmarkPassEnvPtr-8               137547212   8.529 ns/op   0 B/op   0 allocs/op
BenchmarkAccessAllFieldsValue-8     136151521   8.725 ns/op   0 B/op   0 allocs/op
BenchmarkAccessAllFieldsPtr-8       135343456   8.848 ns/op   0 B/op   0 allocs/op
```

---

## Run 2

```bash
BenchmarkSetupEnvValue-8           1000000000   0.3035 ns/op   0 B/op   0 allocs/op
BenchmarkSetupEnvPtr-8             1000000000   0.3032 ns/op   0 B/op   0 allocs/op
BenchmarkPassEnvValue-8              78534236  14.42 ns/op    0 B/op   0 allocs/op
BenchmarkPassEnvPtr-8               140446792   8.762 ns/op   0 B/op   0 allocs/op
BenchmarkAccessAllFieldsValue-8     135499732   8.964 ns/op   0 B/op   0 allocs/op
BenchmarkAccessAllFieldsPtr-8       100000000  13.34 ns/op    0 B/op   0 allocs/op
```

---

## Run 3

```bash
BenchmarkSetupEnvValue-8           1000000000   0.3050 ns/op   0 B/op   0 allocs/op
BenchmarkSetupEnvPtr-8             1000000000   0.3095 ns/op   0 B/op   0 allocs/op
BenchmarkPassEnvValue-8              81528906  14.50 ns/op    0 B/op   0 allocs/op
BenchmarkPassEnvPtr-8               100000000  11.43 ns/op    0 B/op   0 allocs/op
BenchmarkAccessAllFieldsValue-8     127898892   8.413 ns/op   0 B/op   0 allocs/op
BenchmarkAccessAllFieldsPtr-8       140015856   9.473 ns/op   0 B/op   0 allocs/op
```

---

# Result Interpretation

## Environment Creation

Both approaches show nearly identical performance:

* ~0.30 ns/op
* 0 allocations

This happens because the Go compiler optimizes both cases aggressively using:

* inline optimization
* escape analysis
* dead code elimination

This means object creation cost is negligible in this benchmark.

---

## Passing the Struct to Functions

This is where the difference becomes significant.

### Value

```text
~14.4 ns/op to ~14.6 ns/op
```

### Pointer

```text
~8.5 ns/op to ~11.4 ns/op
```

Pointer passing is consistently faster because:

* value copies the entire struct
* pointer copies only an address (8 bytes)

---

## Accessing Fields

Both approaches are very close:

* value and pointer have similar read performance

This shows that after the object is already loaded in memory, field access cost becomes almost identical.

---

# Practical Conclusion

For this project, using a pointer is the best choice:

```go
func run(ctx context.Context, env *config.Env)
```

because:

* `config.Env` is already a medium-sized struct
* it is passed across multiple functions
* copying becomes unnecessary overhead

---

# Recommended Rule

## Use value when:

* struct is very small
* immutable copy is desired

## Use pointer when:

* struct has many fields
* shared configuration is used
* object is passed frequently

---

# Final Decision for This Project

Use:

```go
*config.Env
```
