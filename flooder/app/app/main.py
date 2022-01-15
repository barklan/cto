import pathlib
import random
import time
import sys
from typing import Any, Dict
import asyncio

import fastapi as fa
from fastapi.openapi import utils
import faker

from app.core import custom_logging

fake = faker.Faker()

config_path = pathlib.Path(__file__).with_name("logging_config.json")
log = custom_logging.CustomizeLogger.make_logger(config_path)

app = fa.FastAPI(
    title="Test app",
    debug=True,
)

# https://github.com/tiangolo/fastapi/issues/2750#issuecomment-775526951
# this does serialize exceptions to json but they are still
# duplicated by starlette line by line
sys.stdout.reconfigure(encoding="utf-8", errors="backslashreplace")  # type: ignore


async def catch_exceptions_middleware(request: fa.Request, call_next):  # type: ignore
    try:
        return await call_next(request)
    except Exception:
        request_body = await request.body()
        log.bind(
            method=request.method,
            url=request.url,
            query_params=request.query_params,
            path_params=request.path_params,
            headers=request.headers,
            payload=request_body,
        ).error("Internal server error.")
        # Maybe ditch sentry and send this to dedicated server to special api
        # separate from fluentd logging?
        return fa.Response("Internal server error. Sorry.", status_code=500)


app.middleware("http")(catch_exceptions_middleware)


@app.on_event("startup")
def repeat_random_log():
    while True:
        time.sleep(random.randint(1, 3))
        log.info(fake.text())


@app.get("/info/{id}")
async def infolog(id: str):
    log.info(f"Hey! this is a test log! {id}")
    time.sleep(1)
    return "Hey!"


@app.get("/print/{id}")
async def justprint(id: str):
    print(f"This is print! {id}")
    return "Hey"


@app.get("/error/{id}")
async def errorlog(id: str):
    log.error(f"Hey! this is an er ror log! {id}")
    time.sleep(1)
    return "Hey!"
