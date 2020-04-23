# -*- coding: UTF-8 -*-

import argparse
import os
import xml.etree.ElementTree as ElementTree

'''
格式（全部字段均为 Optional，标记 * 符号的为原型数据库需要）
- 案件类别
- QW 全文（*detail/*title字段）
    - WS 文首
    - DSR 当事人
    - SSJL 诉讼记录
    - AJJBQK 案件基本情况
    - CPFXGC 裁判分析过程
    - PJJC 判决结果
    - WW 文尾
        - ...
        - CUS_FGCY 法官成员
            - FGRYXM 法官人员姓名（*judge字段）
    - LARQ 立案日期
    - WSSFMH 文书是否模糊
    - CUS_SJX 时间线
    - CUS_FJD 附加段
    - CUS_CPWS_CPFXGC 裁判分析过程
    - CUS_CPWS_YGSCD 原告诉称段
    - CUS...
    - FT 法条（*law字段）
'''

walk_rules = [
    # tag   field   only
    ('QW', 'value', True),
    ('FT', 'value', False)
]


class RuleNotCompatibleError(Exception):
    def __init__(self, key, value):
        self.message = 'Parsing error (key={}, value={})'.format(key, value)


parser = argparse.ArgumentParser(description='ElasticJury Data Preprocessor')
parser.add_argument('--path', type=str, default='demo', help='Relative path for data to process')


def add_value(element_map, key, value, only):
    if only and key in element_map:
        raise RuleNotCompatibleError(key, value)
    if only:
        element_map[key] = value
    else:
        if key not in element_map:
            element_map[key] = []
        element_map[key].append(value)


def walk(node, element_map):
    for tag, field, only in walk_rules:
        if tag == node.tag:
            add_value(element_map, tag, node.attrib[field], only)
    for child in node:
        walk(child, element_map)


def process(path):
    print('Processing {}'.format(path))
    element_map = {}
    tree = ElementTree.parse(path)
    walk(tree.getroot(), element_map)
    print(element_map)


def run(path):
    for home, dirs, files in os.walk(path):
        for file in files:
            process(os.path.join(home, file))


if __name__ == '__main__':
    args = parser.parse_args()
    run(args.path)
