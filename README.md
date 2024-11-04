# Paver

## About The Project
Paver is able to run configuration scripts in the required order, taking into account dependencies. Where the DAG is first built and then it is bypassed.
## Options
```sh
--workers - number of workers for parallel execution
--workdir - directory with manifest dirs
```
## Simple example
![Alt text](https://github.com/farit2000/paver/blob/master/assets/simple_example.png)

### workders number = 1
```sh
➜  paver git:(master) ✗ ./bin/paver -workers 1 -workdir examples/simple
Running package: node1
Running script: examples/simple/node1_setup/install.sh
Running package: node2
Task node1 completed with result: [NODE_1 SETUP
]
Running script: examples/simple/node2_setup/install.sh
Running package: node3
Task node2 completed with result: [NODE_2 SETUP
]
Running script: examples/simple/node3_setup/install.sh
Task node3 completed with result: [NODE_3 SETUP
]
Running package: node4
Running script: examples/simple/node4_setup/install.sh
Running package: node5
Running script: examples/simple/node5_setup/install.sh
Task node4 completed with result: [NODE_4 SETUP
]
Running package: node6
Running script: examples/simple/node6_setup/install.sh
Task node5 completed with result: [NODE_5 SETUP
]
Running package: node7
Running script: examples/simple/node7_setup/install.sh
Task node6 completed with result: [NODE_6 SETUP
]
Running package: node8
Running script: examples/simple/node8_setup/install.sh
Task node7 completed with result: [NODE_7 SETUP
]
Task node8 completed with result: [NODE_8 SETUP
]
Time: 17.122071417s
```

### workders number = 2
```sh
➜  paver git:(master) ✗ ./bin/paver -workers 2 -workdir examples/simple
Running package: node1
Running script: examples/simple/node1_setup/install.sh
Running package: node2
Running script: examples/simple/node2_setup/install.sh
Running package: node3
Running script: examples/simple/node3_setup/install.sh
Task node1 completed with result: [NODE_1 SETUP
]
Task node2 completed with result: [NODE_2 SETUP
]
Task node3 completed with result: [NODE_3 SETUP
]
Running package: node4
Running script: examples/simple/node4_setup/install.sh
Running package: node5
Running script: examples/simple/node5_setup/install.sh
Task node5 completed with result: [NODE_5 SETUP
]
Running package: node8
Running script: examples/simple/node8_setup/install.sh
Task node4 completed with result: [NODE_4 SETUP
]
Task node8 completed with result: [NODE_8 SETUP
]
Running package: node7
Running script: examples/simple/node7_setup/install.sh
Running package: node6
Running script: examples/simple/node6_setup/install.sh
Task node6 completed with result: [NODE_6 SETUP
]
Task node7 completed with result: [NODE_7 SETUP
]
Time: 9.058575958s
```

### workders number = 5
```sh
➜  paver git:(master) ✗ ./bin/paver -workers 5 -workdir examples/simple
Running package: node1
Running script: examples/simple/node1_setup/install.sh
Running package: node2
Running package: node3
Running script: examples/simple/node3_setup/install.sh
Running script: examples/simple/node2_setup/install.sh
Task node3 completed with result: [NODE_3 SETUP
]
Running package: node5
Running script: examples/simple/node5_setup/install.sh
Task node1 completed with result: [NODE_1 SETUP
]
Task node2 completed with result: [NODE_2 SETUP
]
Running package: node4
Running script: examples/simple/node4_setup/install.sh
Task node5 completed with result: [NODE_5 SETUP
]
Running package: node8
Running script: examples/simple/node8_setup/install.sh
Task node8 completed with result: [NODE_8 SETUP
]
Task node4 completed with result: [NODE_4 SETUP
]
Running package: node7
Running script: examples/simple/node7_setup/install.sh
Running package: node6
Running script: examples/simple/node6_setup/install.sh
Task node7 completed with result: [NODE_7 SETUP
]
Task node6 completed with result: [NODE_6 SETUP
]
Time: 8.039951625s
```

### if graph has cycle
```sh
➜  paver git:(master) ✗ ./bin/paver -workers 5 -workdir examples/simple
Time: 2.304625ms
panic: error while configure workspace: graph validation failed: cycle detected
 ```

## Real example
Example of a real project with dependencies between scripts and packages. Where installed git, docker, neovim and homebrew on Ubuntu.

![Alt text](https://github.com/farit2000/paver/blob/master/assets/real_example.png)


