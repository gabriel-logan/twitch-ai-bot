# Benchmark Analysis: `config.Env` by Value vs Pointer

## Test Environment
* **CPU**: AMD EPYC 7763 64-Core Processor
* **Architecture**: amd64
* **OS**: Linux
* **Go test target**: `./internal/ws/`

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

| Benchmark | Iterations | ns/op | B/op | allocs/op |
|---|---:|---:|---:|---:|
| BenchmarkSetupEnvValue-4 | 14829405 | 84.58 | 48 | 2 |
| BenchmarkSetupEnvPtr-4 | 16707600 | 71.51 | 48 | 2 |
| BenchmarkPassEnvValue-4 | 79060825 | 14.49 | 0 | 0 |
| BenchmarkPassEnvPtr-4 | 161969701 | 7.613 | 0 | 0 |
| BenchmarkAccessAllFieldsValue-4 | 160521747 | 7.693 | 0 | 0 |
| BenchmarkAccessAllFieldsPtr-4 | 164919740 | 7.302 | 0 | 0 |


---

## Run 2

| Benchmark | Iterations | ns/op | B/op | allocs/op |
|---|---:|---:|---:|---:|
| BenchmarkSetupEnvValue-4 | 15675697 | 76.27 | 48 | 2 |
| BenchmarkSetupEnvPtr-4 | 16842288 | 72.63 | 48 | 2 |
| BenchmarkPassEnvValue-4 | 76673356 | 14.77 | 0 | 0 |
| BenchmarkPassEnvPtr-4 | 158406127 | 7.437 | 0 | 0 |
| BenchmarkAccessAllFieldsValue-4 | 158531172 | 7.543 | 0 | 0 |
| BenchmarkAccessAllFieldsPtr-4 | 148907617 | 7.431 | 0 | 0 |


---

## Run 3

| Benchmark | Iterations | ns/op | B/op | allocs/op |
|---|---:|---:|---:|---:|
| BenchmarkSetupEnvValue-4 | 15772764 | 74.68 | 48 | 2 |
| BenchmarkSetupEnvPtr-4 | 17564322 | 74.14 | 48 | 2 |
| BenchmarkPassEnvValue-4 | 80332974 | 14.66 | 0 | 0 |
| BenchmarkPassEnvPtr-4 | 161178534 | 7.590 | 0 | 0 |
| BenchmarkAccessAllFieldsValue-4 | 158545330 | 7.718 | 0 | 0 |
| BenchmarkAccessAllFieldsPtr-4 | 163030546 | 7.334 | 0 | 0 |


---

## Run 4

| Benchmark | Iterations | ns/op | B/op | allocs/op |
|---|---:|---:|---:|---:|
| BenchmarkSetupEnvValue-4 | 16500852 | 74.34 | 48 | 2 |
| BenchmarkSetupEnvPtr-4 | 16681446 | 70.09 | 48 | 2 |
| BenchmarkPassEnvValue-4 | 75061386 | 14.63 | 0 | 0 |
| BenchmarkPassEnvPtr-4 | 144275503 | 7.263 | 0 | 0 |
| BenchmarkAccessAllFieldsValue-4 | 148677422 | 7.715 | 0 | 0 |
| BenchmarkAccessAllFieldsPtr-4 | 158886720 | 7.476 | 0 | 0 |


---

## Run 5

| Benchmark | Iterations | ns/op | B/op | allocs/op |
|---|---:|---:|---:|---:|
| BenchmarkSetupEnvValue-4 | 16892463 | 74.73 | 48 | 2 |
| BenchmarkSetupEnvPtr-4 | 15951242 | 69.77 | 48 | 2 |
| BenchmarkPassEnvValue-4 | 78677354 | 14.69 | 0 | 0 |
| BenchmarkPassEnvPtr-4 | 156451375 | 7.451 | 0 | 0 |
| BenchmarkAccessAllFieldsValue-4 | 152838811 | 8.094 | 0 | 0 |
| BenchmarkAccessAllFieldsPtr-4 | 161217052 | 7.909 | 0 | 0 |


---

# General Result

## Average of All Runs

| Benchmark | Mean ns/op | B/op | allocs/op |
|---|---:|---:|---:|
| BenchmarkSetupEnvValue | 76.92 | 48 | 2 |
| BenchmarkSetupEnvPtr | 71.63 | 48 | 2 |
| BenchmarkPassEnvValue | 14.65 | 0 | 0 |
| BenchmarkPassEnvPtr | 7.471 | 0 | 0 |
| BenchmarkAccessAllFieldsValue | 7.753 | 0 | 0 |
| BenchmarkAccessAllFieldsPtr | 7.490 | 0 | 0 |



### Mean Performance Difference
* **Environment creation**: pointer is about **6.88% faster** on average.
* **Passing to functions**: pointer is about **49.00% faster** on average.
* **Accessing all fields**: pointer is about **3.38% faster** on average.

---

# Result Interpretation

## Environment Creation

Both approaches remain very close in memory behavior:

* value: **76.92 ns/op** average
* pointer: **71.63 ns/op** average
* both allocate **48 B/op**
* both perform **2 allocs/op**

Pointer creation is slightly faster on this machine, but the difference is still small in absolute terms.

---

## Passing the Struct to Functions

This remains the most relevant difference.

### Value

| Metric | Result |
|---|---|
| Average | 14.65 ns/op |
| Range | 14.49 ns/op to 14.77 ns/op |

### Pointer

| Metric | Result |
|---|---|
| Average | 7.471 ns/op |
| Range | 7.263 ns/op to 7.613 ns/op |

Pointer passing is clearly better here because:

* value copies the entire struct
* pointer copies only an address
* the results were stable across all 5 runs

---

## Accessing Fields

Field access is very close between both approaches on this stronger machine.

### Value

| Metric | Result |
|---|---|
| Average | 7.753 ns/op |
| Range | 7.543 ns/op to 8.094 ns/op |

### Pointer

| Metric | Result |
|---|---|
| Average | 7.490 ns/op |
| Range | 7.302 ns/op to 7.909 ns/op |

Here the difference is small, with pointer still slightly ahead on average.

---

# Practical Conclusion

For this project, using a pointer is still the best choice:

```go
func run(ctx context.Context, env *config.Env)
```

because:

* `config.Env` is a medium-sized struct
* it is passed across multiple functions
* pointer passing avoids unnecessary copies
* the new benchmark results are more stable and consistently favor pointer usage

---

# Recommended Rule

## Use value when:
* the struct is very small
* an immutable copy is desired
* the struct is not passed around frequently

## Use pointer when:
* the struct has many fields
* shared configuration is used
* the object is passed frequently
* you want to avoid repeated copying costs

---

# Final Decision for This Project

Use:

```go
*config.Env
```
