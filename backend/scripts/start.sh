#!/bin/bash
python -u ./python_interface/main.py &
sleep 20
./s0backend
