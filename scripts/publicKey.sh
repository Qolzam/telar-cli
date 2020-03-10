#!/bin/bash

echo -n $(openssl ec -in key -pubout -out key.pub)
