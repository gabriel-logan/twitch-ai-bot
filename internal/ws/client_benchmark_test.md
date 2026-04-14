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

| Benchmark                            | Iterations | ns/op | B/op | allocs/op |
| ------------------------------------ | ---------: | ----: | ---: | --------: |
| BenchmarkSetupEnvValue-4             |   17141466 | 73.15 |   48 |         2 |
| BenchmarkSetupEnvPtr-4               |   17716290 | 69.36 |   48 |         2 |
| BenchmarkPassEnvValue-4              |   74660100 | 14.63 |    0 |         0 |
| BenchmarkPassEnvPtr-4                |  162771505 | 7.409 |    0 |         0 |
| BenchmarkPassEnvValueThroughLayers-4 |   82359506 | 14.65 |    0 |         0 |
| BenchmarkPassEnvPtrThroughLayers-4   |  164469148 | 7.473 |    0 |         0 |
| BenchmarkAccessAllFieldsValue-4      |  163398565 | 8.173 |    0 |         0 |
| BenchmarkAccessAllFieldsPtr-4        |  165766435 | 7.299 |    0 |         0 |

---

## Run 2

| Benchmark                            | Iterations | ns/op | B/op | allocs/op |
| ------------------------------------ | ---------: | ----: | ---: | --------: |
| BenchmarkSetupEnvValue-4             |   17593842 | 70.61 |   48 |         2 |
| BenchmarkSetupEnvPtr-4               |   18061474 | 66.67 |   48 |         2 |
| BenchmarkPassEnvValue-4              |   80613002 | 15.16 |    0 |         0 |
| BenchmarkPassEnvPtr-4                |  159829962 | 7.619 |    0 |         0 |
| BenchmarkPassEnvValueThroughLayers-4 |   80489865 | 14.48 |    0 |         0 |
| BenchmarkPassEnvPtrThroughLayers-4   |  161692410 | 7.399 |    0 |         0 |
| BenchmarkAccessAllFieldsValue-4      |  155821522 | 7.236 |    0 |         0 |
| BenchmarkAccessAllFieldsPtr-4        |  160664203 | 7.199 |    0 |         0 |

---

## Run 3

| Benchmark                            | Iterations | ns/op | B/op | allocs/op |
| ------------------------------------ | ---------: | ----: | ---: | --------: |
| BenchmarkSetupEnvValue-4             |   17468119 | 71.19 |   48 |         2 |
| BenchmarkSetupEnvPtr-4               |   17620556 | 67.64 |   48 |         2 |
| BenchmarkPassEnvValue-4              |   81921260 | 14.60 |    0 |         0 |
| BenchmarkPassEnvPtr-4                |  154580797 | 7.517 |    0 |         0 |
| BenchmarkPassEnvValueThroughLayers-4 |   80708550 | 14.63 |    0 |         0 |
| BenchmarkPassEnvPtrThroughLayers-4   |  167237961 | 7.471 |    0 |         0 |
| BenchmarkAccessAllFieldsValue-4      |  160907088 | 7.804 |    0 |         0 |
| BenchmarkAccessAllFieldsPtr-4        |  166759869 | 7.270 |    0 |         0 |

---

## Run 4

| Benchmark                            | Iterations | ns/op | B/op | allocs/op |
| ------------------------------------ | ---------: | ----: | ---: | --------: |
| BenchmarkSetupEnvValue-4             |   15704492 | 70.70 |   48 |         2 |
| BenchmarkSetupEnvPtr-4               |   17542892 | 72.52 |   48 |         2 |
| BenchmarkPassEnvValue-4              |   79590464 | 14.72 |    0 |         0 |
| BenchmarkPassEnvPtr-4                |  162657099 | 7.258 |    0 |         0 |
| BenchmarkPassEnvValueThroughLayers-4 |   81081048 | 14.76 |    0 |         0 |
| BenchmarkPassEnvPtrThroughLayers-4   |  162833307 | 7.579 |    0 |         0 |
| BenchmarkAccessAllFieldsValue-4      |  167034045 | 7.389 |    0 |         0 |
| BenchmarkAccessAllFieldsPtr-4        |  161059974 | 7.651 |    0 |         0 |

---

## Run 5

| Benchmark                            | Iterations | ns/op | B/op | allocs/op |
| ------------------------------------ | ---------: | ----: | ---: | --------: |
| BenchmarkSetupEnvValue-4             |   17319056 | 70.86 |   48 |         2 |
| BenchmarkSetupEnvPtr-4               |   17820525 | 67.28 |   48 |         2 |
| BenchmarkPassEnvValue-4              |   78237195 | 14.67 |    0 |         0 |
| BenchmarkPassEnvPtr-4                |  167981884 | 7.468 |    0 |         0 |
| BenchmarkPassEnvValueThroughLayers-4 |   82642328 | 14.59 |    0 |         0 |
| BenchmarkPassEnvPtrThroughLayers-4   |  158058201 | 7.644 |    0 |         0 |
| BenchmarkAccessAllFieldsValue-4      |  166051496 | 7.477 |    0 |         0 |
| BenchmarkAccessAllFieldsPtr-4        |  164819450 | 7.538 |    0 |         0 |

---

# General Result

## Average of All Runs

| Benchmark                          | Mean ns/op | B/op | allocs/op |
| ---------------------------------- | ---------: | ---: | --------: |
| BenchmarkSetupEnvValue             |      71.30 |   48 |         2 |
| BenchmarkSetupEnvPtr               |      68.69 |   48 |         2 |
| BenchmarkPassEnvValue              |      14.76 |    0 |         0 |
| BenchmarkPassEnvPtr                |      7.454 |    0 |         0 |
| BenchmarkPassEnvValueThroughLayers |      14.62 |    0 |         0 |
| BenchmarkPassEnvPtrThroughLayers   |      7.513 |    0 |         0 |
| BenchmarkAccessAllFieldsValue      |      7.616 |    0 |         0 |
| BenchmarkAccessAllFieldsPtr        |      7.391 |    0 |         0 |

---

### Mean Performance Difference

* **Environment creation**: pointer is about **3.66% faster** on average.
* **Passing to functions**: pointer is about **49.49% faster** on average.
* **Passing through layers**: pointer is about **48.62% faster** on average.
* **Accessing all fields**: pointer is about **2.95% faster** on average.

---

# Result Interpretation

## Environment Creation

Both approaches remain very close in memory behavior:

* value: **71.30 ns/op** average
* pointer: **68.69 ns/op** average
* both allocate **48 B/op**
* both perform **2 allocs/op**

Pointer creation is slightly faster on this machine, but the difference is still small in absolute terms.

---

## Passing the Struct to Functions

This remains the most relevant difference.

### Value

| Metric  | Result                     |
| ------- | -------------------------- |
| Average | 14.76 ns/op                |
| Range   | 14.60 ns/op to 15.16 ns/op |

### Pointer

| Metric  | Result                     |
| ------- | -------------------------- |
| Average | 7.454 ns/op                |
| Range   | 7.258 ns/op to 7.619 ns/op |

Pointer passing is clearly better here because:

* value copies the entire struct
* pointer copies only an address
* the results were stable across all 5 runs

---

## Accessing Fields

Field access is very close between both approaches on this stronger machine.

### Value

| Metric  | Result                     |
| ------- | -------------------------- |
| Average | 7.616 ns/op                |
| Range   | 7.236 ns/op to 8.173 ns/op |

### Pointer

| Metric  | Result                     |
| ------- | -------------------------- |
| Average | 7.391 ns/op                |
| Range   | 7.199 ns/op to 7.651 ns/op |

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
