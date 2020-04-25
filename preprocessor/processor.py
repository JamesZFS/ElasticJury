# -*- coding: UTF-8 -*-

import os
import xml.etree.ElementTree as ElementTree
import utility

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
    'CASEID',  # ID
    'SFXZ',  # 司法行政（只有 8 条涉及到了这个，根本没用的信息）
    'TAG',  # 司法行政的 TAG
]

special_mapping = [
    ('AL', '案例'),
    ('XZQH_P', '行政区划(省)'),
    ('title', '标题'),
    ('SFSS', '是否上诉'),  # 我猜测的，大量 xml 存在这个标签
]

parsing_error_count = 0


class RuleNotCompatibleError(Exception):
    count = 0  # Static error counter

    def __init__(self, key, value):
        RuleNotCompatibleError.count += 1
        self.message = 'Parsing error (key={}, value={})'.format(key, value)


class TagMatchError(Exception):
    count = 0  # Static error counter

    def __init__(self, tag, mapping, path):
        TagMatchError.count += 1
        self.message = 'Tag {} not found or multiply (mapping[\'{}\']={}, file={})'.format(tag, tag, mapping, path)


def add_value(element_map, key, value, only):
    if only and key in element_map:
        raise RuleNotCompatibleError(key, value)
    if only:
        element_map[key] = value
    else:
        if key not in element_map:
            element_map[key] = []
        element_map[key].append(value)


# Cases
#   1. Has nameCN -> ['key'] = nameCN
#   2. Does not have nameCN -> Find mapping
#       2.1 Not found or multiply -> ['key'] = ''
#       2.2 Only -> ['key'] = mapping[tag]
def attrib_filter(mapping, tag, attrib, path):
    if tag in skip_tags:
        return {}

    new_attrib = {'value': attrib['value'] if 'value' in attrib else ''}
    if 'nameCN' in attrib:
        new_attrib['key'] = attrib['nameCN']
    else:
        if (tag not in mapping) or ((tag in mapping) and len(mapping[tag])) > 1:
            raise TagMatchError(tag, mapping[tag] if tag in mapping else '{}', path)
        else:
            new_attrib['key'] = mapping[tag][0]
    return new_attrib


def walk(node, mapping, element_map, builder, path):
    for tag, field, only in walk_rules:
        if tag == node.tag:
            add_value(element_map, tag, node.attrib[field], only)

    tag = node.tag
    try:
        attrib = attrib_filter(mapping, tag, node.attrib, path)
    except TagMatchError as error:
        attrib = {}
        print('>', error.message)

    builder.start(node.tag, attrib)
    for child in node:
        walk(child, mapping, element_map, builder, path)
    builder.end(tag)


def mark_empty(node):
    empty = (len(node.attrib) < 2)
    for child in node:
        mark_empty(child)
        if not child.attrib['empty']:
            empty = False
    node.attrib['empty'] = empty


def rebuild(node, builder):
    if node.attrib['empty']:
        return

    attrib = node.attrib.copy()
    attrib.pop('empty')
    builder.start(node.tag, attrib)
    for child in node:
        rebuild(child, builder)
    builder.end(node.tag)


def analyze(mapping, path):
    try:
        element_map = {}
        builder = ElementTree.TreeBuilder()
        tree = ElementTree.parse(path)

        walk(tree.getroot(), mapping, element_map, builder, path)

        # Filter for useless nodes
        top = builder.close()
        mark_empty(top)
        builder = ElementTree.TreeBuilder()
        rebuild(top, builder)

        element_map['tree'] = ElementTree.tostring(builder.close())
    except ElementTree.ParseError as error:
        global parsing_error_count
        parsing_error_count += 1
        print('Parsing error in {}: {}'.format(path, error))


def process(mapping, path, db_path):
    print('[Processor] Processing {} ...'.format(path), flush=True)

    for key, value in special_mapping:
        mapping[key] = [value]

    all_xmls = utility.get_all_xml_files(path)
    total = len(all_xmls)
    print('[Processor] {} xmls to process'.format(total), flush=True)

    step = current = 0.05
    for index, file in enumerate(all_xmls):
        analyze(mapping, file)
        if (index + 1) / total >= current:
            print('[Processor] {:.0f}% completed'.format(current * 100), flush=True)
            current += step

    global parsing_error_count
    print('[Processor] Done! ({} parsing error, {} bad rules and {} tag match errors)'
          .format(parsing_error_count, RuleNotCompatibleError.count, TagMatchError.count))
