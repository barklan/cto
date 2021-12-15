#!/bin/bash

uvicorn main:app --host 0.0.0.0 --port 8110 --log-level info
