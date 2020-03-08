#!/bin/bash

./fluentbase-server > /dev/null &

curl --data "set fleet truck3 point 33.4762 -112.10923" localhost:9470

curl localhost:9470/set+fleet+truck3+point+33.4762+-112.10923

