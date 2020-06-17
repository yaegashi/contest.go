# contest.go

## Introduction

contest.go is a simple and naive Go library to help you
with writing and testing solutions for competitive programming problems.

For problems with a fixed answer for each test case,
you can use the standard `go test` command to test your solutions.
You can also use `contest-cli` to prepare your solution template
along with test cases retrieved from the web site.

Supported competitive programming services:
- [AtCoder](https://atcoder.jp)

## Installation

Just run the following to install `contest-cli`:

```console
$ go get github.com/yaegashi/contest.go/cmd/contest-cli
```

## Preparing a solution template with test cases

First, you need to make a Go module for your solutions:

```console
$ git init solutions
Initialized empty Git repository in /home/yaegashi/solutions/.git/
$ cd solutions/
$ go mod init solutions
go: creating new go.mod: module solutions
```

You can prepare a generic solution folder
by running `contest-cli` with a folder name:

```console
$ contest-cli foo
2020/06/18 00:27:30 I: Created foo/main.go
2020/06/18 00:27:30 I: Created foo/main_test.go
2020/06/18 00:27:30 I: Created foo/sample1.in.txt
2020/06/18 00:27:31 I: Created foo/sample1.out.txt
```

You are ready to run `go test` with the `sample1` test case:
```console
$ go test -v ./foo
go: finding module for package github.com/yaegashi/contest.go/tester
go: downloading github.com/yaegashi/contest.go v0.0.1
go: found github.com/yaegashi/contest.go/tester in github.com/yaegashi/contest.go v0.0.1
=== RUN   TestContest
=== RUN   TestContest/0:sample1.in.txt
--- PASS: TestContest (0.00s)
    --- PASS: TestContest/0:sample1.in.txt (0.00s)
PASS
ok      solutions/foo   1.319s
```

You can prepare solution folders for AtCoder contests
by passing an URL of them:
```console
$ contest-cli https://atcoder.jp/contests/abc069
2020/06/18 01:16:01 I: Created abc069/a1/main.go
2020/06/18 01:16:01 I: Created abc069/a1/main_test.go
2020/06/18 01:16:01 I: Created abc069/a1/sample1.in.txt
2020/06/18 01:16:01 I: Created abc069/a1/sample1.out.txt
2020/06/18 01:16:01 I: Created abc069/a1/sample2.in.txt
2020/06/18 01:16:01 I: Created abc069/a1/sample2.out.txt
2020/06/18 01:16:01 I: Created abc069/b1/main.go
2020/06/18 01:16:01 I: Created abc069/b1/main_test.go
2020/06/18 01:16:01 I: Created abc069/b1/sample1.in.txt
2020/06/18 01:16:01 I: Created abc069/b1/sample1.out.txt
2020/06/18 01:16:01 I: Created abc069/b1/sample2.in.txt
2020/06/18 01:16:01 I: Created abc069/b1/sample2.out.txt
2020/06/18 01:16:01 I: Created abc069/b1/sample3.in.txt
2020/06/18 01:16:01 I: Created abc069/b1/sample3.out.txt
2020/06/18 01:16:01 I: Created abc069/c1/main.go
2020/06/18 01:16:01 I: Created abc069/c1/main_test.go
2020/06/18 01:16:01 I: Created abc069/c1/sample1.in.txt
2020/06/18 01:16:01 I: Created abc069/c1/sample1.out.txt
2020/06/18 01:16:01 I: Created abc069/c1/sample2.in.txt
2020/06/18 01:16:01 I: Created abc069/c1/sample2.out.txt
2020/06/18 01:16:01 I: Created abc069/c1/sample3.in.txt
2020/06/18 01:16:01 I: Created abc069/c1/sample3.out.txt
2020/06/18 01:16:01 I: Created abc069/c1/sample4.in.txt
2020/06/18 01:16:01 I: Created abc069/c1/sample4.out.txt
2020/06/18 01:16:01 I: Created abc069/c1/sample5.in.txt
2020/06/18 01:16:01 I: Created abc069/c1/sample5.out.txt
2020/06/18 01:16:01 I: Created abc069/d1/main.go
2020/06/18 01:16:01 I: Created abc069/d1/main_test.go
2020/06/18 01:16:01 I: Created abc069/d1/sample1.in.txt
2020/06/18 01:16:01 I: Created abc069/d1/sample1.out.txt
2020/06/18 01:16:01 I: Created abc069/d1/sample2.in.txt
2020/06/18 01:16:01 I: Created abc069/d1/sample2.out.txt
2020/06/18 01:16:01 I: Created abc069/d1/sample3.in.txt
2020/06/18 01:16:01 I: Created abc069/d1/sample3.out.txt
```

The example above shows it creates 4 solution folder
for [AtCoder Beginner Contest 069](https://atcoder.jp/contests/abc069)
with solution templates and test cases.
Example solutions with them can be found at
[diligence.go repository](https://github.com/yaegashi/diligence.go/tree/master/atcoder/abc069).

For [ABC069 Problem D](https://atcoder.jp/contests/abc069/tasks/arc080_b),
you cannot test your solution with contest.go
because the problem accepts multiple possible answers.
Currently contest.go supports no custom judge. 

```console
$ go test -v .
=== RUN   TestContest
=== RUN   TestContest/0:sample1.in.txt
    TestContest/0:sample1.in.txt: tester.go:93: Wrong answer:
        --- want
        +++ got
        @@ -1,2 +1,2 @@
         1 1
        -2 3
        +3 2
=== RUN   TestContest/1:sample2.in.txt
    TestContest/1:sample2.in.txt: tester.go:93: Wrong answer:
        --- want
        +++ got
        @@ -1,3 +0,0 @@
        -1 4 4 4 3
        -2 5 4 5 3
        -2 5 5 5 3
        @@ -0,0 +1,3 @@
        +1 2 2 3 3
        +4 4 4 4 3
        +5 5 5 5 5
=== RUN   TestContest/2:sample3.in.txt
--- FAIL: TestContest (0.00s)
    --- FAIL: TestContest/0:sample1.in.txt (0.00s)
    --- FAIL: TestContest/1:sample2.in.txt (0.00s)
    --- PASS: TestContest/2:sample3.in.txt (0.00s)
FAIL
FAIL    github.com/yaegashi/diligence.go/atcoder/abc069/d1      0.035s
FAIL
```

## TODO

- Support more competitive programming services
- Custom judge integration for problems with multiple possible answers
- Ability to submit solutions to competitive programming services
