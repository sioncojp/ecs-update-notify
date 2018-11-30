# ecs-update-notify

When updating the task with ecs, let me notify you when the task is replaced at ALB level.

## Usage

```shell
make build

./bin/ecs-update-notify -c config.toml
```

## config.toml

```
interval = 30

[[monitor]]
name = "FooCluster"
aws_profile = "foo"
aws_region = "ap-northeast-1"
incoming_webhook = "https://hooks.slack.com/services/....."

[[monitor]]
name = "BarCluster"
aws_profile = "bar"
aws_region = "ap-northeast-1"
incoming_webhook = "https://hooks.slack.com/services/....."
```

# What would you like to solve with this?

cloudwatch event -> get task replacement for ecs -> lambda

With this solution, the task notifies only creation and deletion.

This is very difficult to understand as people see it.

ecs-update-notify is to notify you only when the new container has replaced at the ALB level.

# License
The MIT License

Copyright Shohei Koyama / sioncojp 

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.