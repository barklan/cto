#! /usr/bin/bash
# set -e

bash ./prestart.sh

if [ "${WITH_JUPYTER:-false}" == 'true' ]; then
  bash -c "jupyter lab -y --log-level=40 --port=8895 --ip=0.0.0.0 --allow-root --NotebookApp.token='' \
    --NotebookApp.password='' --NotebookApp.custom_display_url=http://127.0.0.1:8895/lab/tree/backend/app/temp.ipynb"  > /dev/null 2>&1 &
fi

uvicorn app.main:app --host 0.0.0.0 --port 8000 --log-level info
