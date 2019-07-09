from __future__ import print_function

import argparse
import os
import tensorflow as tf
tf.set_random_seed(228)
from model import Artgan
import redis
import psycopg2
from sqlalchemy import create_engine
from sqlalchemy import Table, Column, String, MetaData, DateTime, Integer


meta = MetaData()
images = Table(
    'images', meta, 
    Column('id', Integer, primary_key = True), 
    Column('created_at', DateTime),
    Column('updated_at', DateTime),
    Column('deleted_at', DateTime),
    Column('name', String),
    Column('style', String),
    Column('status', String),
    Column('uploaded_image_path', String),
    Column('stylized_image_path', String),
)


def tf_handler(uploaded_image_path, save_path, style):
    image_size = os.environ["NN_IMAGE_SIZE"]

    tfconfig = tf.ConfigProto(allow_soft_placement=False)
    tfconfig.gpu_options.allow_growth = True
    with tf.Session(config=tfconfig) as sess:
        model = Artgan(sess, style, image_size)
        ok = model.inference(
            uploaded_image_path, 
            save_path,
            resize_to_original=False,
        )
    return ok



if __name__ == '__main__':
    dbconn = "postgresql+psycopg2://postgres:4H0akRk2Pd77@db/postgres"
    db = create_engine(dbconn)

    rconn = "redis://redis:6379/0"
    r = redis.Redis.from_url(rconn)
    ping = r.ping()

    ps = r.pubsub(ignore_subscribe_messages=True)
    ps.subscribe("imageEvents")

    for message in ps.listen():
        tf.reset_default_graph()

        id = message["data"]

        stmt = images.select().where(images.c.id==id)
        result = db.execute(stmt).fetchone()

        uploaded_image_path = result[images.c.uploaded_image_path]
        image_name = uploaded_image_path.split("/")[-1]
        save_path = "/data/stylized/%s" % image_name

        style = result[images.c.style]

        ok = tf_handler(uploaded_image_path, save_path, style)

        if ok:
            stmt = images.update().where(images.c.id==id).\
                values(status="done", stylized_image_path=save_path)
            db.execute(stmt)

        else:
            stmt = images.update().where(images.c.id==id).\
                values(status="failed")
            db.execute(stmt)