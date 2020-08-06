## What it is
* A script to show the increase rate of each irq event on each CPU core
* This is to help you debug performance issue related to irq
## How to build
* Only Linux is supported
```
go build softirq-rate.go 
./softirq-rate
```

## How to use
```
$ ./irq-rate  -h
Usage of ./irq-rate:
  -interval int
    	refresh interval (default 3)
```

## Example Output
```
Refresh Interval: every 3 seconds

                     CPU0           CPU1           CPU2           CPU3
     BLOCK            0/s            0/s            0/s            0/s
        HI            0/s            0/s            0/s            0/s
   HRTIMER            0/s            0/s            0/s            0/s
  IRQ_POLL            0/s            0/s            0/s            0/s
    NET_RX          229/s          283/s           96/s           95/s
    NET_TX            0/s            1/s            0/s            0/s
       RCU           89/s          112/s           89/s           98/s
     SCHED          148/s          148/s          167/s          184/s
   TASKLET           38/s           41/s           81/s           76/s
     TIMER          118/s          129/s          155/s          188/s
```