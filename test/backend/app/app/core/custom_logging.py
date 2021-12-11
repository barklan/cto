import json
import logging
import sys
import typing as t
from pathlib import Path

from loguru import logger


class InterceptHandler(logging.Handler):
    loglevel_mapping = {
        50: "CRITICAL",
        40: "ERROR",
        30: "WARNING",
        20: "INFO",
        10: "DEBUG",
        0: "NOTSET",
    }

    def emit(self, record: t.Any) -> None:
        try:
            level = logger.level(record.levelname).name
        except AttributeError:
            level = self.loglevel_mapping[record.levelno]

        frame, depth = logging.currentframe(), 2
        while frame.f_code.co_filename == logging.__file__:
            frame = frame.f_back  # type: ignore
            depth += 1

        log = logger.bind(request_id="app")
        log.opt(depth=depth, exception=record.exc_info).log(level, record.getMessage())


class CustomizeLogger:
    @classmethod
    def make_logger(cls, config_path: Path) -> t.Any:

        config = cls.load_logging_config(config_path)
        logging_config = config.get("logger")

        logger = cls.customize_logging(
            "".join([logging_config.get("path"), "/", logging_config.get("filename")]),
            level=logging_config.get("level"),
            retention=logging_config.get("retention"),
            rotation=logging_config.get("rotation"),
            format=logging_config.get("format"),
        )
        return logger

    @classmethod
    def customize_logging(cls, filepath: str, level: str, rotation: str, retention: str, format: str) -> t.Any:

        logger.remove()

        serialize = True
        formatting = "{message}"
        logger.add(
            sys.stdout,
            serialize=serialize,
            colorize=True,
            enqueue=True,
            backtrace=True,
            diagnose=True,
            level=level.upper(),
            format=formatting,
        )
        logger.configure(extra={"payload": None, "user": None})

        logging.getLogger("uvicorn").handlers = [InterceptHandler()]
        logging.getLogger("fastapi").handlers = [InterceptHandler()]
        logging.getLogger("uvicorn.access").handlers = [InterceptHandler()]

        return logger.bind(request_id=None, method=None)

    @classmethod
    def load_logging_config(cls, config_path: Path) -> t.Any:
        config = None
        with open(config_path) as config_file:
            config = json.load(config_file)
        return config
