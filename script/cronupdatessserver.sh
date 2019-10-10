#!/bin/sh

nbssa homeip | awk '$4!=0{print $1}' | xargs -I{} nbssa ssserver update -a {} >> /tmp/ssserver_update.log
