#!/bin/bash
 for FILE in ../lib/core/addon/components/*.js;  do
  component=`eval "echo $FILE | cut -d/ -f6"`; 
   yarn docfy-md $component core
 done