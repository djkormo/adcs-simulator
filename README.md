# adcs-simulator
A simple simulator of Active Directory Certificates Services



[![adcs-simulator pipeline](https://github.com/djkormo/adcs-simulator/actions/workflows/adcs-simulator.yaml/badge.svg)](https://github.com/djkormo/adcs-simulator/actions/workflows/adcs-simulator.yaml)


[![Code scanning - action](https://github.com/djkormo/adcs-simulator/actions/workflows/codeql.yaml/badge.svg)](https://github.com/djkormo/adcs-simulator/actions/workflows/codeql.yaml)

make operator-build

make run

make build

make tag-version

make publish-version

make inspect 

make build tag-version publish-version inspect 

docker run -d -p 8443:8443 djkormo/adcs-simulator:0.0.8  -v ca:/usr/local/adcs-sim