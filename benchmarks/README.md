# Benchmarks

This folder contains focused performance benchmarks for RLAAS hot paths.

## Run all benchmarks in this folder

```bash
go test ./benchmarks -run ^$ -bench . -benchmem
```

## Run a specific benchmark

```bash
go test ./benchmarks -run ^$ -bench BenchmarkMemoryStoreIncrement -benchmem
```

```bash
go test ./benchmarks -run ^$ -bench BenchmarkEvaluateFixedWindow -benchmem
```

## Notes

- `BenchmarkMemoryStore*` covers lock-sharded in-memory counter operations.
- `BenchmarkEvaluateFixedWindow` covers SDK decision evaluation throughput with a static policy store.
