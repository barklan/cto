ARG DOCKER_IMAGE_PREFIX=
FROM ${DOCKER_IMAGE_PREFIX}python:3.9

WORKDIR /app

ENV PYTHONDONTWRITEBYTECODE 1
ENV PYTHONUNBUFFERED 1

RUN apt-get update \
    && apt-get install -y gcc python3-dev ufw graphviz graphviz-dev

# install dependencies
RUN pip install --upgrade pip wheel
COPY ./requirements.txt .
RUN pip install -r ./requirements.txt

# copy project
COPY . .
COPY index.html app.js style.css /app/front/

# run entrypoint.sh
RUN chmod +x /app/entrypoint.sh
ENTRYPOINT ["/app/entrypoint.sh"]
