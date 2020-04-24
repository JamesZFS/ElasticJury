import os
import json
import xml.etree.ElementTree as ElementTree


class ConflictMappingError(Exception):
    count = 0  # Static error counter

    def __init__(self, path, tag, key0, key1):
        ConflictMappingError.count += 1
        self.message = 'Conflict mapping of {} and {} under tag {} (path={})'.format(key0, key1, tag, path)


def add_mapping(tag, attrib, path, mapping):
    if 'nameCN' in attrib:
        key = attrib['nameCN']
        if tag in mapping:
            if mapping[tag] != key:
                raise ConflictMappingError(path, tag, mapping[tag], key)
        else:
            mapping[tag] = key


def walk(node, path, mapping):
    try:
        add_mapping(node.tag, node.attrib, path, mapping)
    except ConflictMappingError as error:
        print('>', error.message)

    for child in node:
        walk(child, path, mapping)


def process_mapping(path, mapping):
    tree = ElementTree.parse(path)
    walk(tree.getroot(), path, mapping)


def get_mapping(path):
    print('[Mapping] Processing mapping in {} ... '.format(path), flush=True)
    mapping = {}

    for home, dirs, files in os.walk(path):
        for file in files:
            if file.endswith('.xml'):
                process_mapping(os.path.join(home, file), mapping)

    print('[Mapping] Done! ({} conflicts found)'.format(ConflictMappingError.count), flush=True)
    return mapping


def dump(mapping, mapping_path):
    print('[Mapping] Dumping mapping to {} ...'.format(mapping_path), flush=True)
    string = json.dumps(mapping)
    with open(mapping_path, 'w') as file:
        file.write(string)
    print('[Mapping] Done!', flush=True)
