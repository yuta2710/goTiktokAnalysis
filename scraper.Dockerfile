FROM python:3.10

WORKDIR /code

COPY ./requirements.txt /code/requirements.txt

RUN pip install --upgrade pip 
RUN pip install --no-cache-dir --upgrade -r /code/requirements.txt
RUN playwright install --with-deps

RUN apt-get update && apt-get install xvfb

COPY ./scripts /code/

CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000", "--reload"]