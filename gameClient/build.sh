#!/usr/bin/env bash
export LD_LIBRARY_PATH=./lsquic/src/liblsquic:$LD_LIBRARY_PATH
gcc main.cpp -I./lsquic/include -L./lsquic/src/liblsquic -llsquic -lstdc++ -lraylib -lGL -lm -lpthread -ldl -lrt -lX11 -o main