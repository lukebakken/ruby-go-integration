### Benchmark Results

#### Ruby

```
$ bundle exec ./ruby-multiget-bench.rb
Key count: 131072
                       user     system      total        real
run 1 (mult):     37.390000   4.810000  42.200000 ( 65.852231)
run 1:            35.250000   4.370000  39.620000 (143.269075)
run 2 (mult):     38.480000   4.680000  43.160000 ( 66.983006)
run 2:            34.940000   4.220000  39.160000 (140.027612)
run 3 (mult):     37.870000   4.660000  42.530000 ( 66.332825)
run 3:            36.920000   4.020000  40.940000 (148.334756)
```

```
$ bundle exec ./ruby-multiget-bench.rb
Key count: 131072
                       user     system      total        real
run 1 (mult):     47.400000   6.600000  54.000000 ( 85.177275)
run 1:            41.250000   5.320000  46.570000 (171.371987)
run 2 (mult):     45.900000   6.490000  52.390000 ( 82.321551)
run 2:            39.440000   5.150000  44.590000 (161.507754)
run 3 (mult):     43.730000   6.070000  49.800000 ( 78.410407)
run 3:            38.370000   4.330000  42.700000 (154.365167)
```

#### Go

```
ok /home/lbakken/Projects/src/ruby-go-integration/go-bench 42.682s
ok /home/lbakken/Projects/src/ruby-go-integration/go-bench 42.152s
ok /home/lbakken/Projects/src/ruby-go-integration/go-bench 42.543s
ok /home/lbakken/Projects/src/ruby-go-integration/go-bench 41.482s
```
