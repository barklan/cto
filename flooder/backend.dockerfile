ARG DOCKER_IMAGE_PREFIX=
FROM ${DOCKER_IMAGE_PREFIX}tiangolo/uvicorn-gunicorn-fastapi:python3.8

ARG BUILDKIT_INLINE_CACHE=1

WORKDIR /app/

# Keeps Python from generating .pyc files in the container
ENV PYTHONDONTWRITEBYTECODE=1

# Turns off buffering for easier container logging
ENV PYTHONUNBUFFERED=1

# Install Poetry
RUN pip install wheel

COPY requirements.txt .

RUN pip install -r requirements.txt

# Copy poetry.lock* in case it doesn't exist in the repo
# COPY ./app/pyproject.toml ./app/poetry.lock* /app/

# Allow installing dev dependencies to run tests
# ARG INSTALL_DEV=false
# RUN bash -c "poetry config virtualenvs.create false && \
#     if [ "${INSTALL_DEV}" == 'true' ] ; \
#     then poetry install --no-root ; \
#     else poetry install --no-root --no-dev ; fi"

# ARG INSTALL_JUPYTER=false
# RUN bash -c "if [ $INSTALL_JUPYTER == 'true' ] ; \
#     then pip install --upgrade pip wheel && apt-get update -y && \
#     curl -LO https://github.com/neovim/neovim/releases/latest/download/nvim.appimage && \
#     chmod u+x nvim.appimage && \
#     ./nvim.appimage --appimage-extract && \
#     mv squashfs-root / && ln -s /squashfs-root/AppRun /usr/bin/nvim && \
#     pip install jupyterlab; fi"

COPY ./app /app
ENV PYTHONPATH=/app

COPY ./app/entrypoint.local.sh /
RUN chmod +x /entrypoint.local.sh
