# Benchmark Analysis: `config.Env` by Value vs Pointer

## Test Environment

* **CPU**: Intel(R) Core(TM) i5-8250U CPU @ 1.60GHz
* **RAM**: 12 GB DDR4 @ 1400 MHz
* **OS**: Linux

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
BenchmarkSetupEnvValue-8          	14388418	        82.72 ns/op	      48 B/op	       2 allocs/op
BenchmarkSetupEnvPtr-8            	14722521	        77.50 ns/op	      48 B/op	       2 allocs/op
BenchmarkPassEnvValue-8           	70258984	        15.55 ns/op	       0 B/op	       0 allocs/op
BenchmarkPassEnvPtr-8             	121914757	        9.983 ns/op	       0 B/op	       0 allocs/op
BenchmarkAccessAllFieldsValue-8   	127733354	        9.552 ns/op	       0 B/op	       0 allocs/op
BenchmarkAccessAllFieldsPtr-8     	121534356	        9.271 ns/op	       0 B/op	       0 allocs/op
```

---

## Run 2

```bash
BenchmarkSetupEnvValue-8          	13033200	        83.89 ns/op	      48 B/op	       2 allocs/op
BenchmarkSetupEnvPtr-8            	11982387	        86.05 ns/op	      48 B/op	       2 allocs/op
BenchmarkPassEnvValue-8           	69109017	        15.62 ns/op	       0 B/op	       0 allocs/op
BenchmarkPassEnvPtr-8             	132078182	        15.91 ns/op	       0 B/op	       0 allocs/op
BenchmarkAccessAllFieldsValue-8   	62719278	        18.92 ns/op	       0 B/op	       0 allocs/op
BenchmarkAccessAllFieldsPtr-8     	122925182	        8.700 ns/op	       0 B/op	       0 allocs/op
```

---

## Run 3

```bash
BenchmarkSetupEnvValue-8          	12201316	        87.27 ns/op	      48 B/op	       2 allocs/op
BenchmarkSetupEnvPtr-8            	13877317	        78.86 ns/op	      48 B/op	       2 allocs/op
BenchmarkPassEnvValue-8           	81014121	        15.77 ns/op	       0 B/op	       0 allocs/op
BenchmarkPassEnvPtr-8             	123778851	        9.019 ns/op	       0 B/op	       0 allocs/op
BenchmarkAccessAllFieldsValue-8   	127421020	        9.456 ns/op	       0 B/op	       0 allocs/op
BenchmarkAccessAllFieldsPtr-8     	95513210	        18.56 ns/op	       0 B/op	       0 allocs/op
```

---

## Run 4

```bash
BenchmarkSetupEnvValue-8          	12525980	        84.20 ns/op	      48 B/op	       2 allocs/op
BenchmarkSetupEnvPtr-8            	12757929	        78.66 ns/op	      48 B/op	       2 allocs/op
BenchmarkPassEnvValue-8           	81415234	        14.56 ns/op	       0 B/op	       0 allocs/op
BenchmarkPassEnvPtr-8             	134989942	        9.252 ns/op	       0 B/op	       0 allocs/op
BenchmarkAccessAllFieldsValue-8   	129291007	        9.308 ns/op	       0 B/op	       0 allocs/op
BenchmarkAccessAllFieldsPtr-8     	132179647	        8.979 ns/op	       0 B/op	       0 allocs/op
```

---

## Run 5

```bash
BenchmarkSetupEnvValue-8          	12188198	        83.61 ns/op	      48 B/op	       2 allocs/op
BenchmarkSetupEnvPtr-8            	12357740	        82.49 ns/op	      48 B/op	       2 allocs/op
BenchmarkPassEnvValue-8           	72068565	        15.38 ns/op	       0 B/op	       0 allocs/op
BenchmarkPassEnvPtr-8             	129901725	        9.123 ns/op	       0 B/op	       0 allocs/op
BenchmarkAccessAllFieldsValue-8   	128307884	        15.57 ns/op	       0 B/op	       0 allocs/op
BenchmarkAccessAllFieldsPtr-8     	139923158	        8.803 ns/op	       0 B/op	       0 allocs/op
```

---

# Result Interpretation

## Environment Creation

Both approaches show nearly identical performance:

* ~77-87 ns/op
* 48 B/op
* 2 allocs/op

The difference is negligible - either approach is fine for object creation.

---

## Passing the Struct to Functions

This is where the difference becomes significant.

### Value

```text
~14.5 ns/op to ~15.8 ns/op
```

### Pointer

```text
~9.0 ns/op to ~15.9 ns/op
```

Pointer passing is generally faster because:

* value copies the entire struct (~84 bytes)
* pointer copies only an address (8 bytes)

---

## Accessing Fields

The results are mixed - both approaches show variable performance:

* value: ~9.3 ns/op to ~18.9 ns/op
* pointer: ~8.7 ns/op to ~18.6 ns/op

Neither is consistently better - the CPU cache and branch predictor heavily influence these results.

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
