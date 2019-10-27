# coding: utf-8
import time
import hashlib
import simplejson
from datetime import datetime, date
import pymysql as connector
from flask import current_app as app
import traceback
import hmac
# import hashlib
import json
import requests
import socket
import re
from ops.models import db

 

def format_model_attribute(attr):
    if isinstance(attr, datetime):
        return datetime.strftime(attr, '%Y-%m-%d %H:%M:%S')
    else:
        return attr


# for datatime type to json
class CJsonEncoder(simplejson.JSONEncoder):

    def default(self, obj):
        if isinstance(obj, datetime):
            return obj.strftime('%Y-%m-%d %H:%M:%S')
        if isinstance(obj, date):
            return obj.strftime('%Y-%m-%d')
        if isinstance(obj, time):
            return obj.strftime('%H:%M:%S')
        else:
            return simplejson.JSONEncoder.default(self, obj)


def get_mysql_connection(server, db):
    try:
        dbconfig = {
            'user': app.config['LOGICAL_USER'],
            'passwd': app.config['LOGICAL_AUTH'],
            'charset': 'utf8',
        }
        host, port = server.split(':')
        port = int(port)
        dbconfig['host'] = host
        dbconfig['port'] = port
        dbconfig['db'] = db
        return connector.connect(**dbconfig)
    except Exception as e:
        print(" get_mysql_connection() error : %s " % traceback.format_exc())
        raise e

 