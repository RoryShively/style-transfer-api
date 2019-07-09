from __future__ import division
from __future__ import print_function

import os
import time
# from glob import glob
import tensorflow as tf
import numpy as np
from collections import namedtuple
# from tqdm import tqdm
import multiprocessing

from module import *
from utils import *
# import prepare_dataset
# import img_augm


class Artgan(object):
    def __init__(self, sess, style, image_size):
        batch_size = 1
        image_size = image_size
        total_steps = int(3e5)
        save_freq = 1000
        lr = 0.0002
        ngf = 32
        ndf = 64
        phase = 'inference'
        path_to_content_dataset = ''
        path_to_art_dataset = ''
        discr_loss_weight = 1.
        transformer_loss_weight = 100.
        feature_loss_weight = 100.

        self.model_name = 'model_%s' % style
        self.root_dir = '/models'
        self.checkpoint_dir = os.path.join(self.root_dir, self.model_name, 'checkpoint')
        self.checkpoint_long_dir = os.path.join(self.root_dir, self.model_name, 'checkpoint_long')
        self.sample_dir = os.path.join(self.root_dir, self.model_name, 'sample')
        self.inference_dir = os.path.join(self.root_dir, self.model_name, 'inference')
        self.logs_dir = os.path.join(self.root_dir, self.model_name, 'logs')

        self.sess = sess
        self.batch_size = batch_size
        self.image_size = image_size

        self.loss = sce_criterion

        self.initial_step = 0

        OPTIONS = namedtuple('OPTIONS',
                             'batch_size image_size \
                              total_steps save_freq lr\
                              gf_dim df_dim \
                              is_training \
                              path_to_content_dataset \
                              path_to_art_dataset \
                              discr_loss_weight transformer_loss_weight feature_loss_weight')
        self.options = OPTIONS._make((
            batch_size, image_size,
            total_steps, save_freq, lr,
            ngf, ndf,
            phase == 'train',
            path_to_content_dataset,
            path_to_art_dataset,
            discr_loss_weight, transformer_loss_weight, feature_loss_weight
        ))

        # Create all the folders for saving the model
        if not os.path.exists(self.root_dir):
            os.makedirs(self.root_dir)
        if not os.path.exists(os.path.join(self.root_dir, self.model_name)):
            os.makedirs(os.path.join(self.root_dir, self.model_name))
        if not os.path.exists(self.checkpoint_dir):
            os.makedirs(self.checkpoint_dir)
        if not os.path.exists(self.checkpoint_long_dir):
            os.makedirs(self.checkpoint_long_dir)
        if not os.path.exists(self.sample_dir):
            os.makedirs(self.sample_dir)
        if not os.path.exists(self.inference_dir):
            os.makedirs(self.inference_dir)

        self._build_model()
        self.saver = tf.train.Saver(max_to_keep=2)
        self.saver_long = tf.train.Saver(max_to_keep=None)

    def _build_model(self):
        # ==================== Define placeholders. ===================== #
        with tf.name_scope('placeholder'):
            self.input_photo = tf.placeholder(dtype=tf.float32,
                                                shape=[self.batch_size, None, None, 3],
                                                name='photo')

        # ===================== Wire the graph. ========================= #
        # Encode input images.
        self.input_photo_features = encoder(image=self.input_photo,
                                            options=self.options,
                                            reuse=False)

        # Decode obtained features.
        self.output_photo = decoder(features=self.input_photo_features,
                                    options=self.options,
                                    reuse=False)

    def inference(self, img_path, save_path, resize_to_original=True,
                  ckpt_nmbr=None):

        init_op = tf.global_variables_initializer()
        self.sess.run(init_op)
        print("Start inference.")

        if self.load(self.checkpoint_dir, ckpt_nmbr):
            print(" [*] Load SUCCESS")
        else:
            if self.load(self.checkpoint_long_dir, ckpt_nmbr):
                print(" [*] Load SUCCESS")
            else:
                print(" [!] Load failed...")
                return False

        # # Create folder to store results.
        # if to_save_dir is None:
        #     to_save_dir = os.path.join(self.root_dir, self.model_name,
        #                                'inference_ckpt%d_sz%d' % (self.initial_step, self.image_size))

        print("Starting")
        try:
            save_dir = "/data/stylized"
            if not os.path.exists(save_dir):
                os.makedirs(save_dir)

            img = scipy.misc.imread(img_path, mode='RGB')
            img_shape = img.shape[:2]
            print("1")

            # Resize the smallest side of the image to the self.image_size
            alpha = float(self.image_size) / float(min(img_shape))
            img = scipy.misc.imresize(img, size=alpha)
            img = np.expand_dims(img, axis=0)

            print("2")

            img = self.sess.run(
                self.output_photo,
                feed_dict={
                    self.input_photo: normalize_arr_of_imgs(img),
                })

            print("3")

            img = img[0]
            img = denormalize_arr_of_imgs(img)

            print("4")
            if resize_to_original:
                img = scipy.misc.imresize(img, size=img_shape)
            else:
                pass

            print("5")
            img_name = os.path.basename(img_path)
            scipy.misc.imsave(save_path, img)

            print("Inference is finished.")

            return True
        
        except Exception as e:
            # return e to store Error in image model
            print("Failed")
            return False

    # def save(self, step, is_long=False):
    #     if not os.path.exists(self.checkpoint_dir):
    #         os.makedirs(self.checkpoint_dir)
    #     if is_long:
    #         self.saver_long.save(self.sess,
    #                              os.path.join(self.checkpoint_long_dir, self.model_name+'_%d.ckpt' % step),
    #                              global_step=step)
    #     else:
    #         self.saver.save(self.sess,
    #                         os.path.join(self.checkpoint_dir, self.model_name + '_%d.ckpt' % step),
    #                         global_step=step)

    def load(self, checkpoint_dir, ckpt_nmbr=None):
        print(" [*] Reading latest checkpoint from folder %s." % (checkpoint_dir))
        ckpt = tf.train.get_checkpoint_state(checkpoint_dir)
        if ckpt and ckpt.model_checkpoint_path:
            ckpt_name = os.path.basename(ckpt.model_checkpoint_path)
            self.initial_step = int(ckpt_name.split("_")[-1].split(".")[0])
            print("Load checkpoint %s. Initial step: %s." % (ckpt_name, self.initial_step))
            self.saver.restore(self.sess, os.path.join(checkpoint_dir, ckpt_name))
            return True
        else:
            return False
