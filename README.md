## Build
```bash
go build -o sigReader sigReader.go utils.go
```

## Run
```bash
# show usage
./sigReader -h
# show process 29966 signal mask
./sigReader -pid 29966
# parse signal mask 7be3c0fe28014a03
./sigReader -mask 7be3c0fe28014a03
```
