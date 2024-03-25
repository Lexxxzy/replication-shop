#!/usr/bin/env bash

export PYDEVD_DISABLE_FILE_VALIDATION=1

pip install nbconvert

cp /home/jovyan/work/notebooks/benchmark_chart.ipynb /tmp/benchmark_chart.ipynb

jupyter nbconvert --execute --to pdf --output-dir=/home/jovyan/work/reports \
  --TemplateExporter.exclude_input=True \
  /tmp/benchmark_chart.ipynb
