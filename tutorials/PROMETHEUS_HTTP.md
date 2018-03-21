# Prometheus HTTP Server Executor


This tutorial covers how to register, expose, and execute tests through an
HTTP server.

Test-Engine comes with an engine server for on demand test execution:

    ```
    $ make build

    go build -o "test-executor" ./commands/test-executor
    go build ./commands/engine-server

    $ ./engine-server -help
      Usage of ./engine-server:
        -metrics string
              Metrics to use: default|prometheus (default "default")
        -testDir string
              Path to serve tests from
    ```

`engine-server` is configurable to submit execution metrics via prometheus.  
Additionally, it needs to be told where to locate tests execute.  In this tutorial
we will run it with tests located on the local filesystem:

    ```
    $ ./engine-server -testDir=$(pwd)/tests -metrics=prometheus 2>&1 | jq .
    {
      "level": "info",
      "msg": "testsDir: \"/vagrant_data/go/src/github.com/dm03514/test-engine/tests\"",
      "time": "2018-03-20T20:55:57-04:00"
    }
    {
      "level": "info",
      "msg": "metrics: \"prometheus\"",
      "time": "2018-03-20T20:55:57-04:00"
    }
    ```

Awesome! We've start the test engine and passed it a directory where to find
tests!  Let's execute a test.  But first, in order to do so, we'll need
test dependencies: 

    - A dummy HTTP REST Stub Server
        
        ```
        $ go run tests/commands/stub-server/*.go
        ```
        
    - elasticmq docker
        
        ```
        $ cd tests/services 
        $ docker-compose -f elasticmq.docker-file.yml down && docker-compose -f elasticmq.docker-file.yml up
        ```
        
- First let's ping the server to verify that it has prometheus metrics:
    ```
    $ curl localhost:8080/metrics
    # HELP go_gc_duration_seconds A summary of the GC invocation durations.
    # TYPE go_gc_duration_seconds summary
    go_gc_duration_seconds{quantile="0"} 0
    go_gc_duration_seconds{quantile="0.25"} 0
    go_gc_duration_seconds{quantile="0.5"} 0
    go_gc_duration_seconds{quantile="0.75"} 0
    go_gc_duration_seconds{quantile="1"} 0
    ```
    
- Now we'll execute a test:

    ```
    $ curl -X POST localhost:8080/execute?test=gstreamer.yml
    ```
    
    For a detailed look at the `gstreamer` test, check out the [`gstreamer` tutorial](GSTREAMER_TUTORIAL.md).
    
- The engine server should show detailed log information about test execution:
    Eventually ending with:
    
    ```
    ...
    {"level":"info","msg":"Subprocess.Evaluate() Subprocess out: \"1\\n\", err: %!s(\u003cnil\u003e)","time":"2018-03-20T21:02:16-04:00"}
    {"component":"engine.Run()","execution_id":"716f6384-6866-49aa-9604-6a7455dd833d","level":"debug","more":true,"msg":"\u003c-resultChan","time":"2018-03-20T21:02:16-04:00"}
    {"adding":"poll_analysis_complete","component":"results","level":"info","msg":"Add()","name":"poll_analysis_complete","time":"2018-03-20T21:02:16-04:00"}
    {"component":"engine.Run()","execution_id":"716f6384-6866-49aa-9604-6a7455dd833d","level":"debug","more":false,"msg":"\u003c-resultChan","time":"2018-03-20T21:02:16-04:00"}
    {"level":"info","msg":"IsLastState(), currState 5 : len(states): 5","time":"2018-03-20T21:02:16-04:00"}
    {"level":"info","msg":"prometheus.Record(6.757363905s, \u003cnil\u003e, errToResult(\"pass\"), \"gstreamer\", \u0026{MetricVec:0xc420062180}","time":"2018-03-20T21:02:16-04:00"}
    {"level":"info","msg":"SUCCESS!","time":"2018-03-20T21:02:16-04:00"
    ```
    
    SUCCESS!!!
    
- Using the promtheus engine server exposes total text execution timings, test result
    pass/fail, and timings for each one of the individual test states!
    
    ```
    $ curl localhost:8080/metrics
    
    state_duration_bucket{result="pass",state_name="send_sqs_message",test_name="gstreamer",le="+Inf"} 1
    state_duration_sum{result="pass",state_name="send_sqs_message",test_name="gstreamer"} 0.613081266
    state_duration_count{result="pass",state_name="send_sqs_message",test_name="gstreamer"} 1
    # HELP test_duration_seconds Duration of a complete test with result:pass|fail
    # TYPE test_duration_seconds histogram
    test_duration_seconds_bucket{result="pass",test_name="gstreamer",le="0.005"} 0
    test_duration_seconds_bucket{result="pass",test_name="gstreamer",le="0.01"} 0
    test_duration_seconds_bucket{result="pass",test_name="gstreamer",le="0.025"} 0
    test_duration_seconds_bucket{result="pass",test_name="gstreamer",le="0.05"} 0
    test_duration_seconds_bucket{result="pass",test_name="gstreamer",le="0.1"} 0
    test_duration_seconds_bucket{result="pass",test_name="gstreamer",le="0.25"} 0
    ```