#!/bin/sh
url="http://localhost:5000"
read -p '[t/s]: ' gr
if [ "$gr" = "t" ]
then
	curl "$url/transactions"
elif [ "$gr" = "s" ]
then
	read sd
	curl -X POST "$url/settle" --data $sd
else
	echo t or s
fi
./$(basename $0) && exit
