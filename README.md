
## Run Tests:
```
cd ./pkg
go test
```

## Run Benchmarks:
```
cd ./pkg
go test -bench . -benchmem -count 10 > 10_runs_bench.txt
```