
go test -bench="Fib$" -cpuprofile=cpu.pprof .
go tool pprof -text cpu.pprof
go tool pprof -png cpu.pprof
go tool pprof -web cpu.pprof