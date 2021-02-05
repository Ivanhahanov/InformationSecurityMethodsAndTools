#!/bin/bash
miss_password_request()
{
  password=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9'  | head -c 10);
  curl -X POST --header "X-Forwarded-For: $ip" -F "username=admin" -F "password=$password" http://localhost:8080/login --silent --output /dev/null;
}
hacker()
{
  ip=$(($RANDOM % 255)).$(($RANDOM % 255)).$(($RANDOM % 255)).$(($RANDOM % 255));
  for i in {1..10}
  do
  password=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9'  | head -c 10);
  curl -X POST --header "X-Forwarded-For: $ip" -F "username=admin" -F "password=$password" http://localhost:8080/login --silent --output /dev/null;
  sleep 0.5
  done
}
loop_num=100
hackers_attack=$(($RANDOM % $loop_num))
for i in $(seq 0 $loop_num)
do
  if [ $i -eq $hackers_attack ]
  then
  hacker
  else
  ip=$(($RANDOM % 255)).$(($RANDOM % 255)).$(($RANDOM % 255)).$(($RANDOM % 255));
  for ((j=0; j<= $(($RANDOM % 2 - 1)); j++));
  do
    miss_password_request;
    sleep 1;
  done
  curl -X POST --header "X-Forwarded-For: $ip" -F 'username=admin' -F 'password=admin' http://localhost:8080/login --silent --output /dev/null;
  sleep 2;
  fi
done
