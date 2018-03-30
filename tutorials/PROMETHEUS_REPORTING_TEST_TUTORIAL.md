# Creating and Executing Prometheus Reported Test

This tutorial will simulate developing a complex test, executing it, and getting
prometheus metrics from it, for this we will assume there is an API to develop against.
This tutorial will cover:
    - test definition
    - test CLI execution
    - test incremental development
    - test engine loggign overview
    - test http/subprocess
    - state variable overrides/substitutions
    - test CLI execution
    - Prometheus HTTP Test server configuration
    - Prometheus HTTP Test server execution
    
For this tutorial let's pretend we have 3 separate microservice systems and we'd like
to test a complete workflow through them.  The systems are:

    - Account Service (HTTP/REST)
    - Analysis Service (REACTIVE/SQS)
    - Result Service (HTTP/REST)
    
The account service manages accounts for analysis, the anlysis service performs some sort
of analysis and then provides a result, and the result service exposes analysis results.
This tutorial will provision a test account through the account service, fire off an SQS
message for analysis, and then poll the result service until completion, or timeout.  After
writing a test to do this we will learn how to register this test through an HTTP server
so that we can execute it at will, or on an interval.

    
- Start stub server in background
This tutorial includes a stub server with all endpoints created:

    ```
    $ go run tests/commands/stub-server/*.go
    ```

- Create a test template, and the first state that will create a new account per test run
    
    ```
    # tests/gstreamer.yml
    ---
    name: gstreamer
    states:
      - name: create_fake_account
        fulfillment:
          type: noop.Noop
        action:
          type: http.Http
          url: http://localhost:9999/gstreamer/account/create
          method: POST
        transition_conditions:
          - type: assertions.IntEqual
            using_property: status_code
            to_equal: 200
    ```
    
- Let's execute our test including the first state!

    ```
    $ go run commands/test-executor/main.go -test $(pwd)/tests/gstreamer.yml
    
    # or with piping log files to stdout and jq 
    # $  go run commands/test-executor/main.go -test $(pwd)/tests/gstreamer.yml 2>&1 | jq .
    ```
    
    ```javascript
    {
      "component": "Fileloader.Load()",
      "filename": "gstreamer.yml",
      "level": "info",
      "msg": "loading_test",
      "path": "/vagrant_data/go/src/github.com/dm03514/test-engine/tests/gstreamer.yml",
      "time": "2018-03-20T08:48:33-04:00"
    }    
    ...
    {
      "against": 201,
      "component": "assertions.IntEqual",
      "execution_id": "8700c13f-ede3-48a3-a87d-f148b8d5c148",
      "level": "info",
      "msg": "Evaluate()",
      "time": "2018-03-20T08:48:33-04:00",
      "to_equal": 200,
      "using_property": "status_code"
    }
  
  {
    "component": "engine.Run()",
    "execution_id": "8700c13f-ede3-48a3-a87d-f148b8d5c148",
    "level": "debug",
    "more": true,
    "msg": "<-resultChan",
    "time": "2018-03-20T08:48:33-04:00"
  }
  {
    "level": "panic",
    "msg": "200 != 201",
    "time": "2018-03-20T08:48:33-04:00"
  }
    ```
    
    Our first test execution!!!!
   
    Test-Engine logs are focused on verbosity and auditability.  Since tests and
    states can be long lived there needs to be a way to tie log statements
    with potentially large intervals between statements, together.
    
    TODO - access execution_id from inside of tests 
    TODO - execution_id on top level panic
    
- Let's update the assertion to be correct and rerun the test.  Since we're creating a new resource so the return code
should be 201 instead of 200:
    
    ```
        transition_conditions:
          - type: assertions.IntEqual
            using_property: status_code
            to_equal: 201
    ```
    
    ```
    $ go run commands/test-executor/main.go -test $(pwd)/tests/gstreamer.yml 
    
    {
      "component": "test-executor.main",
      "level": "info",
      "msg": "SUCCESS",
      "time": "2018-03-20T08:58:08-04:00"
    }
    ```
    
    SUCCESS!!!!!!!!!!!!!!

- Now that we've created a new account to use for this test we need to extract
the account id that was just created.  In order to do this, Test-Engine allows
for bash to be executed so that more specific tools can be used.  The next 
state drops into bash to manipulate the API response from the previous state using `jq`

    ```
     - name: parse_account_id
        fulfillment:
          type: noop.Noop
        action:
          type: shell.Subprocess
          command_name: sh
          args:
            - "-c"
            - "printf '$ACCOUNT_CREATE_REST_RESPONSE' | jq -r .AccountID | tr -d '\n'"
          overrides:
            - from_state: create_fake_account
              using_property: body
              to_replace: $ACCOUNT_CREATE_REST_RESPONSE
        transition_conditions:
          - type: assertions.IntEqual
            using_property: returncode
            to_equal: 0
    ```
    
    - There above leverages leverages the `Subprocess` `action` to allow the use of
        bash
    - Since this is just bash it allows for the trivial testing of commands, so that
        individual bash commands don't have to be tested in the context of the Test-Engine
        framework, allowing for faster feedback loops of development.
        
        ie. By taking the API rest response from the first step we can develop
            the manipulation we want locally in our shell:
        
        ```
        $ printf '{\"AccountID\":\"ID-CREATED\"}\n' | jq .AccountID | tr -d '\n'
        
        ID-CREATED
        ```
    - Now that the action to manipulate the API response has been validated
    we now need to inject it into the action.  This is performed by declaring an `override`
    which references the response from the first state:
    
        ```
            - "printf '$ACCOUNT_CREATE_REST_RESPONSE' | jq -r .AccountID | tr -d '\n'"
          overrides:
            - from_state: create_fake_account
              using_property: body
              to_replace: $ACCOUNT_CREATE_REST_RESPONSE
        ```
        
        - `from_state` is the state `name` that we will get a result from
        - `using_property` is the property of the state result that we will inject into the action
            in this case we will be using the `http.Http` action `body` result.
        - `to_replace` is the variable in the action that will be substituted
        
    - Running the test with this new state shows:
    
        ```
        $ go run commands/test-executor/main.go -test $(pwd)/tests/gstreamer.yml
       
        {
          "args": [
            "-c",
            "printf '{\"AccountID\":\"ID-CREATED\"}\n' | jq -r .AccountID | tr -d '\n'"
          ],
          "command": "sh",
          "component": "shell.Subprocess",
          "error": null,
          "execution_id": "eab38f79-6556-4f92-8266-c83df5c0c9bb",
          "level": "info",
          "msg": "CombinedOutput()",
          "output": "ID-CREATED",
          "time": "2018-03-20T11:17:36-04:00"
        } 
       
        {
          "component": "test-executor.main",
          "level": "info",
          "msg": "SUCCESS",
          "time": "2018-03-20T11:17:36-04:00"
        } 
        ```
       
        WIN! 
        
- To recap we've created a test which makes an API request to create a new resource, and
have used bash and `jq` to manipulate that response into a form we can use in subsequent states.
The next step is using the id of the created account in a command to generate a test SQS, then publishing
that message, and finally, polling for the results of that message.

- Building the SQS message shares a number of similarities with the other steps:

    ```
      - name: build_sqs_message
        fulfillment:
          type: noop.Noop
        action:
          type: shell.Subprocess
          command_name: sh
          args:
            - "-c"
            - >
              printf '{"account_id":"$ACCOUNT_ID","test":"$UUID_test_msg_identifier"}'
          overrides:
            - from_state: parse_account_id
              using_property: output
              to_replace: $ACCOUNT_ID
        transition_conditions:
          - type: assertions.IntEqual
            using_property: returncode
            to_equal: 0
    ```

    - The above simulates deferring to a command inside of the repo.  One can imagine
        that the data enqueued on SQS needs to be encoding using protobuf or some
        other serialization format, but also needs unique dynamic data from
        the current test execution Executing the task shows:
        
    ```
    $ go run commands/test-executor/main.go -test $(pwd)/tests/gstreamer.yml
    
    {
      "args": [
        "-c",
        "printf '{\"account_id\":\"ID-CREATED\",\"test\":\"75726509-98f2-446a-b73e-0b327dc9c30e\"}'\n"
      ],
      "command": "sh",
      "component": "shell.Subprocess",
      "execution_id": "feb969c0-95bc-4329-bc1c-9b20e4d334ff",
      "level": "info",
      "msg": "Execute()",
      "time": "2018-03-20T13:32:34-04:00"
    }
    {                                                                                                                                                                    [47/2974]
      "args": [
        "-c",
        "printf '{\"account_id\":\"ID-CREATED\",\"test\":\"75726509-98f2-446a-b73e-0b327dc9c30e\"}'\n"
      ],
      "command": "sh",
      "component": "shell.Subprocess",
      "error": null,
      "execution_id": "feb969c0-95bc-4329-bc1c-9b20e4d334ff",
      "level": "info",
      "msg": "CombinedOutput()",
      "output": "{\"account_id\":\"ID-CREATED\",\"test\":\"75726509-98f2-446a-b73e-0b327dc9c30e\"}",
      "time": "2018-03-20T13:32:34-04:00"
    }
    
    {
      "component": "test-executor.main",
      "level": "info",
      "msg": "SUCCESS",
      "time": "2018-03-20T13:32:34-04:00"
    }
    ```
    
    - `$UUID_{{ VARNAME }}` is a special template variable that Test-Engine will replace
        with a UUID.  Each `VARNAME` will receive the same UUID value.
        
    - Additionally, this state substitutes the `parse_account_id` state output (the ACCOUNT_ID)
        and injects it into the command
        
        ```
              printf '{"account_id":"$ACCOUNT_ID","test":"$UUID_test_msg_identifier"}'
          overrides:
            - from_state: parse_account_id
              using_property: output
              to_replace: $ACCOUNT_ID
        ```
        
- Now it's time to submit our message to SQS to kick off the analysis 

    - Start elasticmq for local sqs testing
    
        ```
        $ cd tests/services
        $ docker-compose -f elasticmq.docker-file.yml down && docker-compose -f elasticmq.docker-file.yml up
        ```
        
    - Send the previously generated message to sqs using aws cli tools   
        Hopefully this is starting to look familiar:
    
        ```
          - name: send_sqs_message
            fulfillment:
              type: noop.Noop
            action:
              type: shell.Subprocess
              command_name: sh
              args:
                - "-c"
                - >
                  aws --endpoint-url http://localhost:9324 sqs send-message \
                      --queue-url http://localhost:9324/queue/test_queue \
                      --message-attributes '{ "Content-Transfer-Encoding":{ "DataType":"String","StringValue":"base64" }}' \
                      --message-body "`echo '$BASE64_MESSAGE_BODY' | base64 | tr -d '\n'`"
              overrides:
                - from_state: build_sqs_message
                  using_property: output
                  to_replace: $BASE64_MESSAGE_BODY
            transition_conditions:
              - type: assertions.IntEqual
                using_property: returncode
                to_equal: 0
        ```
        
        We're using the awscli tools to send a message to the local aws.
        The body is being encoded to base64 before being sent.
        
        Once again, running this should show us a message is successfully sent:
        
        ```
        {                                                                                                                                                                        [50/6885]
          "args": [
            "-c",
            "aws --endpoint-url http://localhost:9324 sqs send-message \\\n    --queue-url http://localhost:9324/queue/test_queue \\\n    --message-attributes '{ \"Content-Transfer-Encod
        ing\":{ \"DataType\":\"String\",\"StringValue\":\"base64\" }}' \\\n    --message-body \"`echo '{\"account_id\":\"ID-CREATED\",\"test\":\"9414f2b4-fe9d-4160-8fe4-a89ae399669d\"}'
        | base64 | tr -d '\\n'`\"\n"
          ],
          "command": "sh",
          "component": "shell.Subprocess",
          "error": null,
          "execution_id": "94818888-166c-4c54-b90d-1b851a45078f",
          "level": "info",
          "msg": "CombinedOutput()",
          "output": "{\n    \"MD5OfMessageBody\": \"d3316d1cd375b8afc2397b1ec7ffe213\", \n    \"MD5OfMessageAttributes\": \"b5115b0a9e71841323013bac5993d834\", \n    \"MessageId\": \"2fa
        89c31-a43e-457c-8f2f-9f16b8b092cc\"\n}\n",
          "time": "2018-03-20T15:49:38-04:00"
        } 
        
        {
          "component": "test-executor.main",
          "level": "info",
          "msg": "SUCCESS",
          "time": "2018-03-20T15:40:19-04:00"
        }
        ```
        
        We can additionally verify this message has been sent by receiving a message from
        the local queue!!!
        
        ```
        $ aws --endpoint-url http://localhost:9324 sqs receive-message \
            --queue-url http://localhost:9324/queue/test_queue | jq -r '.Messages[0].Body' | base64 -d
            
          {"account_id":"ID-CREATED","test":"9414f2b4-fe9d-4160-8fe4-a89ae399669d"
        ```
        
        WE'VE SENT THE DYNAMICALLY CREATED MESSAGE!

- Up until this point, the whole test has been provisioning the test data and
sending off the message for processing.  Now we will focus on polling until 
we detect the message has been processed.  Let's pretend processing ends
at a restful api, and we will use the ACCOUNT_ID to correlate the end result with 
our test.

- We'll poll the API until processing has completed:

    ```
      - name: poll_analysis_complete
        fulfillment:
          type: poll.Poller
          interval: 100ms
          timeout: 5m
        action:
          type: http.Http
          url: http://localhost:9999/gstreamer/analysis_complete
        transition_conditions:
          - type: assertions.IntEqual
            using_property: status_code
            to_equal: 200
          - type: assertions.Subprocess
            command_name: sh
            args:
              - "-c"
              - "echo '$body' | jq -e '.results | length'"
            using_property: body
            to_equal: "1\n"
    ``` 
    
    This introduces a couple of new concepts.  THe first is the `poll.Poller` fulfillment
    strategy.  This fulfiller will try the action at the given interval until the
    state transition conditions have passed or until the timeout is reached.
    
    The next new component is the `assertions.Subprocess` which allows for shell commands
    to be executed as the assertion.  In this case we are taking the `body` of the action
    result and using `jq` to manipulate it, and get the `length` of the results.
    
    This state is considered fulfilled when there is a single `result` and the status
    code of the http call is `200`.

    If you execute this task, you'll see it poll for a little until the results are 
    returned, and `SUCCESS` is reached:
    
        {
          "against": 204,
          "component": "assertions.IntEqual",
          "execution_id": "a26c5f2b-7e87-47c6-aaf3-5ef7e8ea0eec",
          "level": "info",
          "msg": "Evaluate()",
          "time": "2018-03-20T20:17:34-04:00",
          "to_equal": 200,
          "using_property": "status_code"
        }
       
- YAY! WOOT we now have a declaritive language agnostic test!!!!!!!!!!!! WHat do we do 
now with it?? [Check out how to expose the test through an HTTP server, which
allows for metric reporting on your test.](PROMETHEUS_HTTP.md)
