#!/usr/bin/env bash

mygeth=$1
dd=$2

echo "mygeth: ${mygeth}"
echo "dd: ${dd}"

networkid=666
datadir=./${dd}/
port=303${dd}
rpcport=81${dd}
verbosity=6
nodekey=./${dd}/${dd}.key
logdir=./${dd}.log
account=`$mygeth --datadir=$datadir account list|head -n1|perl -ne '/([a-f0-9]{40})/ && print $1'`
password=./${dd}/pwd.txt
bootnodes="\"enode://c89d3fd5a361fc54203e19d6bc5ff2472a97f46ac4dfe17609f49a9af6fe1fe8d97b998dc3d4638de220a1677fde48484cd02b36d1cba157ce029cd31ff3d871@127.0.0.1:30301\""

if [ "$dd" = "01" ]; then
    echo "01 is the bootnode"
    command="$mygeth --datadir  $datadir --networkid $networkid --ipcdisable --rpcapi="db,eth,net,web3,personal,web3,exch" --rpc --port $port --rpcport $rpcport  -verbosity $verbosity --nodekey $nodekey  --unlock $account --password $password console 2>> $logdir"
else
    echo "$dd is not the bootnode"
    command="$mygeth --datadir  $datadir --networkid $networkid --ipcdisable --rpcapi="db,eth,net,web3,personal,web3,exch" --rpc --port $port --rpcport $rpcport  -verbosity $verbosity --nodekey $nodekey  --unlock $account --password $password --bootnodes $bootnodes console 2>> $logdir"
fi

echo $command

eval $command