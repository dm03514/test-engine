# Test-engine
An asynchronous test toolkit powering end to end and integration level tests.

Test-engine provides:
- Centralized framework 
- Primitives for common async operations (polling, event/queue subscriptions, timeouts, etc)
- Test result and duration reporting
- Unified logging format
- Flexible support for extension in bash or go

Test engine is an extensible asynchronous test execution framework.  It is aimed at supporting end-to-end and integration testing.
Test-engine allows for modeling of tests as a series of simple states using a delcaritive yaml syntax.
Test-engine then schedules, executes, applies state transitions, and reports on test statuses.  Test-engine is a great use for:

- Language agnostic tests (think [travis-ci](https://travis-ci.org/) of testing)
- Simple declaritave test configuration (also supporting yaml), suited for both programmers and non-programmers
- Ability to drop down and configure tests directly in code 
- Extensible shared test framework: A go core allows for leveraging of rich go client library ecosystem
- Easily extensible component-based design
  - Teams can easily extend test-engine by adding their own libraries or primitivies

Writing tests in Test-engine is amazingly fast after learning the test executing flow and toolkit.  It enables reliable standardized and uniform tests which can be written in minutes. It allows tests to be written in the same framework regardless of service language.  Test-engine is a powerful test framework which benefits from a composable, extensible, component based approach.  Effort to extending the framework is minimal and then accessible by all service teams independent of service language.

Test-engine minimizes the logic in end-to-end tests, by declaring states and conditions necessary to transition each state.  Test-engine then takes care of transitioning your tests through those states and the concurrency, allowing you to focus on what your service does, and not the test infrastructure.

## Problem
Writing reliable functional/acceptance (e2e) tests is hard.  They are notorious for being flaky, time consuming, and difficult 
to maintain because they require:

- multiple protocols
- multiple systems
- waiting (timing)
- asynchronous systems
- complex test logic

Test-engine aims to separate the task definition, from how it is executed, and how each step of the task is transitioned.  
It should allow for complex multi-step functional/acceptance (e2e) tests, often involving concurrent operations, to be statically 
defined as a list of state definitions.


## An Asynchronous Test Toolkit

Test engine is a toolkit of primitives which enable easy reliable testing of events 
which happen in the future.  Testing through public interfaces often necessitates asynchrounous
interactions.  Consider an image analysis service exposed through an http interface:

Post an image and service returns a resource ID
resource then contains the status

Consider asynchronous queue based microservices.  Services read from an input queue, processes a 
message, and either 

The test glue code to apply input then wait or poll for output is often boilerplate and due to the nature of tests often not 
rigorously directly tested or reused.


Suppose we're testing a service with mysql as a dependency.  We're using docker to manage
mysql locally and in CI.  Before tests are executed mysql needs to be started:

```
$ docker run -e MYSQL_ALLOW_EMPTY_PASSWORD=1 -p 3306:2206 --name=mysql mysql
$ ./run-service-test-that-depends-on-mysql
```

The issue is that the tests should only start after mysql is fully initialized and is able to 
accept and process connections.  If not this could result in a race condition where
the test suite is launched and some tests fail due to mysql not being able to handle them
but later tests execute fine.  Having a test suite and mysql communicate across processes
makes it difficult to synchronize.  

A common first approach to solving this problem is timeout.  This approach is suboptimal and 
should be considered unsynchronized.  

```
$ docker run -e MYSQL_ALLOW_EMPTY_PASSWORD=1 -p 3306:2206 --name=mysql mysql
$ sleep 10
$ ./run-service-test-that-depends-on-mysql
```

This test starts mysql and pauses to give mysql an opportunity to fully initialze.  If the 
timeout vaule is too long than time is wasted.  If the timeout chosen is too short it can 
be flaky because mysql is still not fully initialized.  Additionally, just because
a timeout is sufficient locally does not mean its sufficient across build servers.
A CI server or a coworkers server could have significantly different specs and the timeout
may not provide a sufficient amount of time for mysql to initialize, resulting in flaky tests.

A much better approach is to actively check to detect when mysql is ready and only yield
after it is proven to be ready.  This is the approach that wait-for-it takes. The issue is
that wait for it is an extremely narrowly focused tool, it doesn't provide many knobs
for configuration or building test libraries.  Test-engine is like wait-for-it for
all protocols, test assertions, parameterizable test templates, test execution metrics, and is 
customizable.

Suppose from manual observation and experience mysql returning a cursor in the CLI is
the most reliable way to determine if mysql is available to accept connections and 
handle our tests:

```
$ mysql -h 127.0.0.1 -u root -e "SELECT 1"
+---+
| 1 |
+---+
| 1 |
+---+
```

With mysql unable to accept connections and process queries the command returns:

```
$ mysql -h 127.0.0.1 -u root -e "SELECT version"
ERROR 2003 (HY000): Can't connect to MySQL server on '127.0.0.1' (111
```

Using Test-engine we can create and schedule a check when mysql is ready:

```yaml
name: mysql_ready
states:
  - name: mysql_ready
    fulfillment:
      type: poll.Poller
      interval: 500ms
      timeout: 10s
    action:
      type: shell.Subprocess
      command_name: sh
      args:
        - "-c"
        - "mysql -h 127.0.0.1 -u root -e \"SELECT 1\""
    transition_conditions:
      - type: assertions.IntEqual
        using_property: returncode
        to_equal: 0
      - type: assertions.StringEqual
        using_property: output
        to_equal: "1\n1\n"
```

At the core of the check above is a single state `mysql_ready`.  Each state has 3 different
components: `action`, `transition_conditions` and `fulfillment`.  The `action`.  This is the active check we will execute
to determine if mysql is ready or not.  Next are the `transition_conditions` required for 
the check to complete.  In this case the test will pass when the returncode (exit status) is
succesful (`0`) AND the output is equal to:

```
1
1
```

Finally, the `fulfillment` strategy is how the action is scheduled.  In this case
the action is being scheduled using polling fulfillment.  Thi will execute the action
every 500ms until the timeout is reached.

Executing the test without mysql running shows:

```
$ ./test-executor -test=$(pwd)/tests/mysql-initialized.yml

{
  "component": "Fileloader.Load()",
  "filename": "mysql-initialized.yml",
  "level": "info",
  "msg": "loading_test",
  "path": "/vagrant_data/go/src/github.com/dm03514/test-engine/tests/mysql-initialized.yml",
  "time": "2018-03-13T15:20:49-04:00"
}
{
  "component": "NewFromYaml()",
  "level": "debug",
  "msg": "parsing_test",
  "test_name": "mysql_ready",
  "time": "2018-03-13T15:20:49-04:00"
}
...
{
  "component": "Engine.ExecuteState()",
  "current_state_index": 0,
  "level": "info",
  "msg": "executing",
  "state": "mysql_ready",
  "time": "2018-03-13T15:20:49-04:00"
}
{
  "component": "poll.Poller",
  "interval": "500ms",
  "level": "info",
  "msg": "starting_poller",
  "time": "2018-03-13T15:20:49-04:00",
  "timeout": "10s"
}
...
{
  "args": [
    "-c",
    "mysql -h 127.0.0.1 -u root -e \"SELECT 1\""
  ],
  "command": "sh",
  "component": "shell.Subprocess",
  "error": "exit status 1",
  "level": "info",
  "msg": "CombinedOutput()",
  "output": "ERROR 2003 (HY000): Can't connect to MySQL server on '127.0.0.1' (111)\n",
  "time": "2018-03-13T15:20:50-04:00"
}
{
  "against": 1,
  "component": "assertions.IntEqual",
  "level": "info",
  "msg": "Evaluate()",
  "time": "2018-03-13T15:20:50-04:00",
  "to_equal": 0,
  "using_property": "returncode"
}
...
{
  "level": "panic",
  "msg": "Timeout \"10s\" exceeded",
  "time": "2018-03-13T15:20:59-04:00"
}
```

Since mysql isn't running the action is executed until the timeout is reached in which 
case the command exists with exit stats 1.

Running this test after starting mysql should allow for a reliable way to detect when mysql
is initialized and tests are ready:

```
$ docker run -e MYSQL_ALLOW_EMPTY_PASSWORD=1 -p 3306:3306 mysql                         [171/1309]
  Initializing database
  2018-03-13T19:42:05.326400Z 0 [Warning] TIMESTAMP with implicit DEFAULT value is deprecated. Please use --explicit_defaults_for_timestamp server option (see documentation
   for more details).
  2018-03-13T19:42:05.471859Z 0 [Warning] InnoDB: New log files created, LSN=45790
  2018-03-13T19:42:05.500569Z 0 [Warning] InnoDB: Creating foreign key constraint system tables.
  2018-03-13T19:42:05.556332Z 0 [Warning] No existing UUID has been found, so we assume that this is the first time that this server has been started. Generating a new UUID
  : 9bfb87d8-26f6-11e8-888a-0242ac110002.
  2018-03-13T19:42:05.557756Z 0 [Warning] Gtid table is not ready to be used. Table 'mysql.gtid_executed' cannot be opened.
  2018-03-13T19:42:05.558795Z 1 [Warning] root@localhost is created with an empty password ! Please consider switching off the --initialize-insecure option.
  2018-03-13T19:42:05.748067Z 1 [Warning] 'user' entry 'root@localhost' ignored in --skip-name-resolve mode.
  ...
  2018-03-13T19:42:12.438733Z 0 [Warning] 'db' entry 'sys mysql.sys@localhost' ignored in --skip-name-resolve mode.
  2018-03-13T19:42:12.438745Z 0 [Warning] 'proxies_priv' entry '@ root@localhost' ignored in --skip-name-resolve mode.
  2018-03-13T19:42:12.441160Z 0 [Warning] 'tables_priv' entry 'user mysql.session@localhost' ignored in --skip-name-resolve mode.
  2018-03-13T19:42:12.441187Z 0 [Warning] 'tables_priv' entry 'sys_config mysql.sys@localhost' ignored in --skip-name-resolve mode.
  2018-03-13T19:42:12.450256Z 0 [Note] Event Scheduler: Loaded 0 events
  2018-03-13T19:42:12.450488Z 0 [Note] mysqld: ready for connections.
  Version: '5.7.21'  socket: '/var/run/mysqld/mysqld.sock'  port: 3306  MySQL Community Server (GPL)
```

```
$ ./test-executor -test=$(pwd)/tests/mysql-initialized.yml

...
{
  "component": "engine.Run()",
  "execution_id": "826b3931-6621-4edb-b666-b55ccd1bee8b",
  "level": "debug",
  "more": false,
  "msg": "<-resultChan",
  "time": "2018-03-13T15:44:55-04:00"
}
{
  "level": "info",
  "msg": "IsLastState(), currState 1 : len(states): 1",
  "time": "2018-03-13T15:44:55-04:00"
}
{
  "component": "test-executor.main",
  "level": "info",
  "msg": "SUCCESS",
  "time": "2018-03-13T15:44:55-04:00"
}
```

SUCCESS!!  The mysq-initialized test now polls until mysql docker server is ready for
connections!!!!

