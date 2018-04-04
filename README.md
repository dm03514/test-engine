# Test-Engine

An asynchronous test toolkit powering end to end and integration level tests.

[![Go Report Card](https://goreportcard.com/badge/github.com/dm03514/test-engine)](https://goreportcard.com/report/github.com/dm03514/test-engine)

Test-Engine provides:
- Centralized framework 
- Primitives for common async operations (polling, event/queue subscriptions, timeouts, etc)
- Test result and duration reporting
- Unified logging format
- Flexible support for extension in bash or go

Test-Engine is an extensible asynchronous test execution framework.  It is aimed at supporting end-to-end and integration testing.
Test-Engine allows for modeling of tests as a series of simple states using a declarative yaml syntax.
Test-Engine then schedules, executes, applies state transitions, and reports on test statuses.  Test-Engine is a great use for:

- Language agnostic tests (think [travis-ci](https://travis-ci.org/) of testing)
- Simple declarative test configuration (also supporting yaml), suited for both programmers and non-programmers
- Ability to drop down and configure tests directly in code 
- Extensible shared test framework: A go core allows for leveraging of rich go client library ecosystem
- Easily extensible component-based design
  - Teams can easily extend Test-Engine by adding their own libraries or primitivies

Writing tests in Test-Engine is amazingly fast after learning the test executing flow and toolkit.  It enables reliable standardized and uniform tests which can be written in minutes. 
It allows tests to be written in the same framework regardless of service language.  
Test-Engine is a powerful test framework which benefits from a composable, extensible, component based approach.  
Effort to extending the framework is minimal and then accessible by all service teams independent of service language.

Test-Engine minimizes the logic in end-to-end tests, by declaring states and conditions necessary to transition each state.  Test-Engine then takes care of transitioning your tests through those states and the concurrency, allowing you to focus on what your service does, and not the test infrastructure.

## Problem
Writing reliable functional/acceptance (e2e) tests is hard.  They are notorious for being flaky, time consuming, and difficult 
to maintain because they require:

- multiple protocols
- multiple systems
- waiting (timing)
- asynchronous systems
- complex test logic

Test-Engine aims to separate the task definition, from how it is executed, and how each step of the task is transitioned.  
It should allow for complex multi-step functional/acceptance (e2e) tests, often involving concurrent operations, to be statically 
defined as a list of state definitions.


## Additional Resources

- [Test-Engine Introduction](https://medium.com/dm03514-tech-blog/introducing-test-engine-an-asynchronous-test-toolkit-5ca0883a0f4b)
    - Service Level Tests
    - System Level Tests
    - Service Level Objective Tests
- Creating Tests
    - [Creating and Executing Service Test Using CLI & YAML](tutorials/creating_and_executiong_service_test_using_the_cli.md)
        This tutorial covers everything involved with designing, developing, debugging and executing a non-trivial service test.
    - Create and Executing Custom Tests
- Components/Concepts
    - Action
    - Fulfillment
    - Transition Conditions
    - Overrides
- Deploying Tests 
    - CLI YAML Service Level Integration
    - CLI Custom Test Integration
    - HTTP Service/System Level Integration
    - Prometheus/SLO integration
        - [Creating, Exposing, and Executing Test to Gather SLO Metrics](tutorials/SLO_TEST_EXPOSED_THROUGH_PROMETHEUS_HTTP.md)
        This tutorial covers registering a created test with the HTTP Executor, so that test execution is exposed over REST and reports execution metrics through prometheus.
    - Custom Tests HTTP
- Debugging Tests
- Custom Components
