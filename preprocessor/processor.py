# -*- coding: UTF-8 -*-

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

# Some fixed rules
walk_rules = [
    # tag, field, only
    ('QW', 'value', True),
    ('FT', 'value', False),
    ('FGRYXM', 'value', False),
]

skip_tags = [
    'writ',  # 文章开始
    'FT',  # 法条
    'FTRY',  # 法条冗余
    'CASEID', # ID
]


class RuleNotCompatibleError(Exception):
    count = 0  # Static error counter

    def __init__(self, key, value):
        RuleNotCompatibleError.count += 1
        self.message = 'Parsing error (key={}, value={})'.format(key, value)


class TagNotFoundError(Exception):
    count = 0  # Static error counter

    def __init__(self, tag):
        TagNotFoundError.count += 1
        self.message = 'Tag {} not found'.format(tag)


def add_value(element_map, key, value, only):
    if only and key in element_map:
        raise RuleNotCompatibleError(key, value)
    if only:
        element_map[key] = value
    else:
        if key not in element_map:
            element_map[key] = []
        element_map[key].append(value)


def attrib_filter(mapping, tag, attrib):
    if tag in skip_tags:
        return {}

    new_attrib = {}
    if (tag in mapping) and ('value' in attrib):
        new_attrib['value'] = attrib['value']
    elif tag not in mapping:
        raise TagNotFoundError(tag)
    return new_attrib


def walk(node, mapping, element_map, builder):
    for tag, field, only in walk_rules:
        if tag == node.tag:
            add_value(element_map, tag, node.attrib[field], only)

    tag = node.tag
    try:
        attrib = attrib_filter(mapping, tag, node.attrib)
    except TagNotFoundError as error:
        attrib = {}
        print('>', error.message)

    builder.start(node.tag, attrib)
    for child in node:
        walk(child, mapping, element_map, builder)
    builder.end(tag)


def analyze(mapping, path):
    element_map = {}
    builder = ElementTree.TreeBuilder()
    tree = ElementTree.parse(path)

    walk(tree.getroot(), mapping, element_map, builder)

    element_map['tree'] = ElementTree.tostring(builder.close())


def process(mapping, path, db_path):
    print('[Processor] Processing {} ...'.format(path), flush=True)

    for home, dirs, files in os.walk(path):
        for file in files:
            if file.endswith('.xml'):
                analyze(mapping, os.path.join(home, file))

    print('[Processor] Done! ({} bad rules and {} tags not found)'
          .format(RuleNotCompatibleError.count, TagNotFoundError.count))
