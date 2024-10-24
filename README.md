# Paver

## About The Project
Paver is able to run configuration scripts in the required order, taking into account dependencies. Where the DAG is first built and then it is bypassed.
## Example
```sh
➜  paver git:(master) ✗ ./bin/paver examples 
Running package: git
Running script: examples/git_setup/install.sh
Task git completed with result: [GIT SETUP
]
Running package: nvim
Running script: examples/nvim_setup/install.sh
Running package: docker
Running script: examples/docker_setup/install.sh
Running script: examples/nvim_setup/configure.sh
Task docker completed with result: [DOCKER SETUP
]
Task nvim completed with result: [NVIM SETUP
 ]
 ```