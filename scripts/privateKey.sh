#!/bin/bash

echo -n $(openssl ecparam -genkey -name prime256v1 -noout -out key)
