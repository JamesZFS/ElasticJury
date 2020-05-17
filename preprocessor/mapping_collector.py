# -*- coding: UTF-8 -*-

import json
import xml.etree.ElementTree as ElementTree

from utility import *

parsing_error_count = 0


class ConflictMappingError(Exception):
    count = 0  # Static error counter

    def __init__(self, path, tag, key0, key1):
        ConflictMappingError.count += 1
        self.message = 'Conflict mapping of {} and {} under tag {} (path={})'.format(key0, key1, tag, path)


def add_mapping(tag, attrib, path, mapping):
    if 'nameCN' in attrib:
        key = attrib['nameCN']
        if tag not in mapping:
            mapping[tag] = []
        if key not in mapping[tag]:
            mapping[tag].append(key)


def walk(node, path, mapping):
    try:
        add_mapping(node.tag, node.attrib, path, mapping)
    except ConflictMappingError as error:
        log_info('Error', error.message)

    for child in node:
        walk(child, path, mapping)


def process_mapping(path, mapping):
    try:
        tree = ElementTree.parse(path)
        walk(tree.getroot(), path, mapping)
    except ElementTree.ParseError as error:
        global parsing_error_count
        parsing_error_count += 1
        log_info('Error', 'Parsing error in {}: {}'.format(path, error))


def get_mapping(path):
    log_info('Mapping', 'Processing mapping in {} ... '.format(path))
    mapping = {}

    all_xmls = get_all_xml_files(path)
    total = len(all_xmls)
    log_info('Mapping', '{} xmls to process'.format(total))

    step = current = 0.05
    for index, file in enumerate(all_xmls):
        process_mapping(file, mapping)
        if (index + 1) / total >= current:
            log_info('Mapping', '{:.0f}% completed'.format(current * 100))
            current += step

    global parsing_error_count
    log_info('Mapping', 'Done! ({} parsing error)'.format(parsing_error_count))
    return mapping


def dump(mapping, mapping_path):
    log_info('Mapping', 'Dumping mapping to {} ...'.format(mapping_path))
    string = json.dumps(mapping)
    with open(mapping_path, 'w') as file:
        file.write(string)
    log_info('Mapping', 'Done')
