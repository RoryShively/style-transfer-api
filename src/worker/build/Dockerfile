FROM python:2.7

ENV PYTHONUNBUFFERED=0
ENV CUDA_VISIBLE_DEVICES=0

# RUN     apt-get update && \
#     apt-get install -y \
#     postgresql-client \
#     redis-tools

WORKDIR /app

COPY requirements.txt /app
RUN pip install -r requirements.txt

COPY src/ /app

CMD ["python", "./main.py"]