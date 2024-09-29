FROM python:3

WORKDIR /jwt

COPY app.py /jwt/app.py
COPY templates /jwt/templates
COPY static /jwt/static


RUN apt-get update  \
    && useradd -m ctf-player \
    && chown -R ctf-player /jwt

USER ctf-player

RUN pip install --no-cache-dir Flask pyjwt

EXPOSE 4054

CMD ["python3", "app.py"]
