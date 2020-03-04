#!/usr/bin/env bash

#wget https://fluentbase.org/wp-content/uploads/2019/04/fluentbase-binary.tar.gz

git clone https://github.com/gvsafronov/fluentbase/
cd fluentbase
make
mv tile38-cli fluentbase-cli
mv tile38-server fluentbase-server
mv tile38-benchmark fluentbase-benchmark
mv tile38-luamemtest fluentbase-luamemtest
chmod +x fluentbase-server && chmod +x fluentbase-cli
 ./fluentbase-server 

#wget https://fluentbase.org/wp-content/uploads/2019/04/flunix-binary.tar.gz
#wget https://github.com/gvsafronov/flunix/
#    tar xzvf flunix-binary.tar.gz
#        ./flunix



