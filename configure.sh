#!/usr/bin/env bash

#wget https://fluidb.org/wp-content/uploads/2019/04/fluidb-binary.tar.gz

git clone https://github.com/gvsafronov/fluidb/
cd fluidb
git clone https://github.com/gvsafronov/flunix/
make
mv tile38-cli clif
mv tile38-server fluidb-serv
mv tile38-benchmark fluidb-benchmark
mv tile38-luamemtest fluidb-luamemtest
chmod +x fluidb-server && chmod +x clif
 ./flunix

#wget https://fluidb.org/wp-content/uploads/2019/04/flunix-binary.tar.gz
#wget https://github.com/gvsafronov/flunix/
#    tar xzvf flunix-binary.tar.gz
#        ./flunix



