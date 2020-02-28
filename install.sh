#!/usr/bin/env bash

#wget https://fluentbase.org/wp-content/uploads/2019/04/fluentbase-binary.tar.gz
wget https://github.com/gvsafronov/fluentbase/

    tar xzvf fluentbase-binary.tar.gz
        cd fluentbase-binary
            chmod +x fluentbase-server && chmod +x fluentbase-cli
 ./fluentbase-server &

#wget https://fluentbase.org/wp-content/uploads/2019/04/flunix-binary.tar.gz
wget https://github.com/gvsafronov/flunix/
    tar xzvf flunix-binary.tar.gz
        ./flunix



