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
```bash
BenchmarkSetupEnvValue-4                14829405                84.58 ns/op           48 B/op          2 allocs/op
BenchmarkSetupEnvPtr-4                  16707600                71.51 ns/op           48 B/op          2 allocs/op
BenchmarkPassEnvValue-4                 79060825                14.49 ns/op            0 B/op          0 allocs/op
BenchmarkPassEnvPtr-4                   161969701                7.613 ns/op           0 B/op          0 allocs/op
BenchmarkAccessAllFieldsValue-4         160521747                7.693 ns/op           0 B/op          0 allocs/op
BenchmarkAccessAllFieldsPtr-4           164919740                7.302 ns/op           0 B/op          0 allocs/op
```

---

## Run 2
```bash
BenchmarkSetupEnvValue-4                15675697                76.27 ns/op           48 B/op          2 allocs/op
BenchmarkSetupEnvPtr-4                  16842288                72.63 ns/op           48 B/op          2 allocs/op
BenchmarkPassEnvValue-4                 76673356                14.77 ns/op            0 B/op          0 allocs/op
BenchmarkPassEnvPtr-4                   158406127                7.437 ns/op           0 B/op          0 allocs/op
BenchmarkAccessAllFieldsValue-4         158531172                7.543 ns/op           0 B/op          0 allocs/op
BenchmarkAccessAllFieldsPtr-4           148907617                7.431 ns/op           0 B/op          0 allocs/op
```

---

## Run 3
```bash
BenchmarkSetupEnvValue-4                15772764                74.68 ns/op           48 B/op          2 allocs/op
BenchmarkSetupEnvPtr-4                  17564322                74.14 ns/op           48 B/op          2 allocs/op
BenchmarkPassEnvValue-4                 80332974                14.66 ns/op            0 B/op          0 allocs/op
BenchmarkPassEnvPtr-4                   161178534                7.590 ns/op           0 B/op          0 allocs/op
BenchmarkAccessAllFieldsValue-4         158545330                7.718 ns/op           0 B/op          0 allocs/op
BenchmarkAccessAllFieldsPtr-4           163030546                7.334 ns/op           0 B/op          0 allocs/op
```

---

## Run 4
```bash
BenchmarkSetupEnvValue-4                16500852                74.34 ns/op           48 B/op          2 allocs/op
BenchmarkSetupEnvPtr-4                  16681446                70.09 ns/op           48 B/op          2 allocs/op
BenchmarkPassEnvValue-4                 75061386                14.63 ns/op            0 B/op          0 allocs/op
BenchmarkPassEnvPtr-4                   144275503                7.263 ns/op           0 B/op          0 allocs/op
BenchmarkAccessAllFieldsValue-4         148677422                7.715 ns/op           0 B/op          0 allocs/op
BenchmarkAccessAllFieldsPtr-4           158886720                7.476 ns/op           0 B/op          0 allocs/op
```

---

## Run 5
```bash
BenchmarkSetupEnvValue-4                16892463                74.73 ns/op           48 B/op          2 allocs/op
BenchmarkSetupEnvPtr-4                  15951242                69.77 ns/op           48 B/op          2 allocs/op
BenchmarkPassEnvValue-4                 78677354                14.69 ns/op            0 B/op          0 allocs/op
BenchmarkPassEnvPtr-4                   156451375                7.451 ns/op           0 B/op          0 allocs/op
BenchmarkAccessAllFieldsValue-4         152838811                8.094 ns/op           0 B/op          0 allocs/op
BenchmarkAccessAllFieldsPtr-4           161217052                7.909 ns/op           0 B/op          0 allocs/op
```

---

# General Result

## Average of All Runs
```text
BenchmarkSetupEnvValue     mean: 76.92 ns/op   48 B/op   2 allocs/op
BenchmarkSetupEnvPtr       mean: 71.63 ns/op   48 B/op   2 allocs/op
BenchmarkPassEnvValue      mean: 14.65 ns/op   0 B/op    0 allocs/op
BenchmarkPassEnvPtr        mean: 7.471 ns/op  0 B/op    0 allocs/op
BenchmarkAccessAllFieldsValue mean: 7.753 ns/op  0 B/op    0 allocs/op
BenchmarkAccessAllFieldsPtr   mean: 7.490 ns/op  0 B/op    0 allocs/op
```


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
```text
average: 14.65 ns/op
range: 14.49 ns/op to 14.77 ns/op
```

### Pointer
```text
average: 7.471 ns/op
range: 7.263 ns/op to 7.613 ns/op
```

Pointer passing is clearly better here because:

* value copies the entire struct
* pointer copies only an address
* the results were stable across all 5 runs

---

## Accessing Fields

Field access is very close between both approaches on this stronger machine.

### Value
```text
average: 7.753 ns/op
range: 7.543 ns/op to 8.094 ns/op
```

### Pointer
```text
average: 7.490 ns/op
range: 7.302 ns/op to 7.909 ns/op
```

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
