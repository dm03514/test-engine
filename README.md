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


## Additional Resources
- [Test-Engine Introduction](https://medium.com/dm03514-tech-blog/introducing-test-engine-an-asynchronous-test-toolkit-5ca0883a0f4b)

