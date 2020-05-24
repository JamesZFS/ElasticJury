# -*- coding: UTF-8 -*-

import MySQLdb
import jieba
import jieba.analyse
import re
import xml.etree.ElementTree as ElementTree

from database import MySQLWrapper
from entries import *
from utility import *

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
    # tag, field, only, used in tag_field
    ('QW', 'value', True, False),  # 全文
    ('FT', 'value', False, False),  # 法条
    ('FGRYXM', 'value', False, False),  # 法官人员姓名
    ('AJLB', 'value', False, True),  # 案件类别
    ('AJLX', 'value', False, True),  # 案件类型
    ('AJSJ', 'value', False, True),  # 案件涉及
    ('WSMC', 'value', False, True),  # 文书名称
    ('SPCX', 'value', False, True),  # 审判程序
    ('AY', 'value', False, True),  # 案由
    ('SJLY_AY', 'value', False, True),  # 涉及领域_案由
    ('WZAY', 'value', False, True),  # 完整案由
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


def add_value(entry, key, value, only):
    if only and key in entry:
        raise RuleNotCompatibleError(key, value)
    if only:
        entry[key] = value
    else:
        if key not in entry:
            entry[key] = []
        entry[key].append(value)


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


def walk(node, mapping, entry, builder, path):
    for tag, field, only, _ in walk_rules:
        if tag == node.tag and field in node.attrib:
            add_value(entry, tag, node.attrib[field], only)

    tag = node.tag
    try:
        attrib = attrib_filter(mapping, tag, node.attrib, path)
    except TagMatchError as error:
        attrib = {}
        print('>', error.message)

    builder.start(node.tag, attrib)
    for child in node:
        walk(child, mapping, entry, builder, path)
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
    entry = {}
    try:
        builder = ElementTree.TreeBuilder()
        tree = ElementTree.parse(path)

        walk(tree.getroot(), mapping, entry, builder, path)

        # Filter for useless nodes
        top = builder.close()
        mark_empty(top)
        builder = ElementTree.TreeBuilder()
        rebuild(top, builder)

        top = builder.close()
        entry['tree'] = ElementTree.tostring(top, encoding='unicode')
    except ElementTree.ParseError as error:
        global parsing_error_count
        parsing_error_count += 1
        log_info('Error', 'Parsing error in {}: {}'.format(path, error))
    return entry


def reduce_count_weights(items):
    counter = {}
    for item in items:
        item = item.lower()
        value = counter.get(item, 0)
        counter[item] = value + 1
    total = len(items)
    return [(k, v / total) for k, v in counter.items()]


def reduce_words(items, idf_dict):
    counter = {}
    total = 0
    for item in items:
        item = item.lower()
        value = counter.get(item, 0)
        counter[item] = value + 1
        total += 1
    return [(k, v / total * idf_dict[k]) for k, v in counter.items()]


def collect_entries(items_with_weights, table_name, key_name, case_id):
    entries = []
    for item, weight in items_with_weights:
        entry = IndexEntry(table_name, key_name, item, case_id, weight)
        entries.append(entry)
    return entries


law_pattern = re.compile('《.+》')


def reduce_laws(laws):
    reduced = {}
    for k, v in reduce_count_weights(laws):
        law = law_pattern.findall(k)
        if len(law) == 1:
            law = law[0]
            reduced[law] = v + reduced.get(law, 0)
        else:
            log_info('Debug', 'Failed to parse law {}'.format(k))
        reduced[k] = v + reduced.get(k, 0)
    return [(k, v) for k, v in reduced.items()]


stopwords = set([line.strip() for line in open('stopwords.txt', 'r', encoding='utf-8').readlines()])


def insert_into_database(database, idf_dict, entry):
    if len(entry) == 0:
        return

    detail = entry.get('QW', '')
    tree = entry.get('tree', '')

    xml_tags = []
    for tag, _, only, used_in_tag in walk_rules:
        if used_in_tag:
            assert not only
            for to_split in entry.get(tag, []):
                xml_tags.extend(filter(lambda x: len(x) > 0, re.split('[ 、]', to_split)))
    xml_tags = [(k, 1) for k in filter(lambda x: len(x) < 32, set(xml_tags))]
    judges = list(filter(lambda x: len(x) < 16, entry.get('FGRYXM', [])))

    arrays = [
        (reduce_words(filter(
            lambda w: (w not in stopwords) and (len(w.strip()) > 0)
                      and len(w) < 16, jieba.lcut(detail)), idf_dict),
         'WordIndex', 'word'),
        [reduce_count_weights(judges), 'JudgeIndex', 'judge'],
        (reduce_laws(entry.get('FT', [])), 'LawIndex', 'law'),
        # AllowPos 是词性位置
        (xml_tags, 'TagIndex', 'tag')
    ]

    # Insert case
    extract = lambda x, i: [item[0] for item in x[i][0]]
    case = CaseEntry(extract(arrays, 1), extract(arrays, 2), extract(arrays, 3), detail, tree)
    case_id = database.insert(case)

    # Insert index
    for array in arrays:
        entries = collect_entries(array[0], array[1], array[2], case_id)
        database.insert_many(entries)

    log_info('ShowCase', '')
    log_info('ShowCase', 'case_id={}'.format(case_id))
    log_info('ShowCase', 'detail={}'.format(detail))
    log_info('ShowCase', 'tree={}'.format(tree))
    words = arrays[0][0]
    words.sort(key=lambda x: x[1], reverse=True)
    log_info('ShowCase', 'words={}'.format(words))
    log_info('ShowCase', 'judges={}'.format(arrays[1][0]))
    log_info('ShowCase', 'laws={}'.format(arrays[2][0]))
    log_info('ShowCase', 'tags={}'.format(arrays[3][0]))


def process(mapping, idf_dict, path, db_password, drop):
    log_info('Jieba', 'Initializing jieba ...')
    jieba.initialize()
    log_info('Processor', 'Processing {} ...'.format(path))

    for key, value in special_mapping:
        mapping[key] = [value]

    all_xmls = get_all_xml_files(path)
    total = len(all_xmls)
    log_info('Processor', '{} xmls to process'.format(total))

    database = MySQLWrapper(drop=drop, password=db_password)

    step = current = 0.005
    database_error_count = 0
    for index, file in enumerate(all_xmls):
        log_info('Processor', 'Processing {}'.format(file))
        entry = analyze(mapping, file)
        insert_into_database(database, idf_dict, entry)
        # try:
        #     insert_into_database(database, idf_dict, entry)
        # except MySQLdb.DatabaseError as error:
        #     database_error_count += 1
        #     log_info('Debug', 'Insert error: {} at file {}'.format(error, file))
        if (index + 1) / total >= current:
            log_info('Processor', '{:.2f}% completed'.format(current * 100))
            current += step

    global parsing_error_count
    log_info('Processor', '({} parsing error, {} database error, {} bad rules and {} tag match errors)'
             .format(parsing_error_count, database_error_count, RuleNotCompatibleError.count, TagMatchError.count))
    database.close()
