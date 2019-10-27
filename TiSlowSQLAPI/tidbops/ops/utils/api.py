# -*- coding: UTF-8 -*-
from flask import jsonify

# define statu_dics here
RS200 = {'code': 200, 'msg': 'OK all right.'}
RS201 = {'code': 201, 'msg': 'All created.'}
RS204 = {'code': 204, 'msg': 'All deleted.'}
RS400 = {'code': 400, 'msg': 'Bad request.'}
RS401 = {'code': 401, 'msg': 'Authorize failed'}
RS403 = {'code': 403, 'msg': 'You can not do this.'}
RS404 = {'code': 404, 'msg': 'No result matched.'}
RS500 = {'code': 500, 'msg': 'Internal error.'}
RS520 = {'code': 520, 'msg': 'Internal exception msg'}
RSYES = {'code': 1, 'msg': 'Successful.'}
RSNO = {'code': 0, 'msg': 'Fail'}


def fullResponse(statu_dic, data):
    return jsonify({'status': statu_dic, 'data': data})


def statusResponse(statu_dic):
    return jsonify({'status': statu_dic})
