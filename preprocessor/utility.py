# -*- coding: UTF-8 -*-

import os


class UnimplementedError(Exception):
    def __init__(self):
        self.message = 'Unimplemented code'


def unimplemented(_obj=None):
    raise UnimplementedError()


def get_all_xml_files(path):
    xmls = []
    for home, dirs, files in os.walk(path):
        for file in files:
            if file.endswith('.xml'):
                xmls.append(os.path.join(home, file))
    return xmls


def log_info(tag, info, flush=True, end='\n'):
    print('[{}] {}'.format(tag, info), flush=flush, end=end)


def log_exit(tag, info, flush=True, end='\n'):
    print('[{}] {}'.format(tag, info), flush=flush, end=end)
    exit(1)
