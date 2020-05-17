import jieba
import json
import math
import xml.etree.ElementTree as ElementTree

from utility import *

parsing_error_count = 0


def parse_words(text, idf_dict):
    for word in set(jieba.lcut(text)):
        count = idf_dict.get(word, 0)
        idf_dict[word] = count + 1


def walk(node, counter_dict):
    if node.tag == 'QW':
        parse_words(node.attrib['value'], counter_dict)
    for child in node:
        walk(child, counter_dict)


def process_idf_dict(path, counter_dict):
    try:
        tree = ElementTree.parse(path)
        walk(tree.getroot(), counter_dict)
    except ElementTree.ParseError:
        global parsing_error_count
        parsing_error_count += 1
        log_info('Error', 'Parsing error in {}'.format(path))


def get_idf_dict(path):
    log_info('Jieba', 'Initializing jieba ...')
    jieba.initialize()
    log_info('IdfDict', 'Processing idf dict in {} ... '.format(path))
    counter_dict = {}

    all_xmls = get_all_xml_files(path)
    total = len(all_xmls)
    log_info('IdfDict', '{} xmls to process'.format(total))

    step = current = 0.01
    for index, file in enumerate(all_xmls):
        process_idf_dict(file, counter_dict)
        if (index + 1) / total >= current:
            log_info('IdfDict', '{:.0f}% completed'.format(current * 100))
            current += step

    log_info('IdfDict', 'Running idf computing for words')
    idf_dict = {}
    for k, v in counter_dict.items():
        idf_dict[k] = math.log(total / v)

    global parsing_error_count
    log_info('IdfDict', 'Done! ({} parsing error)'.format(parsing_error_count))
    return idf_dict


def dump(idf_dict, idf_dict_path):
    log_info('IdfDict', 'Dumping idf dict to {} ...'.format(idf_dict_path))
    string = json.dumps(idf_dict)
    with open(idf_dict_path, 'w') as file:
        file.write(string)
    log_info('IdfDict', 'Done')
