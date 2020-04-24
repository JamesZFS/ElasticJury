import json
import xml.etree.ElementTree as ElementTree
import utility


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

        # One-to-one mapping
        # if tag in mapping:
        #     if mapping[tag] != key:
        #         raise ConflictMappingError(path, tag, mapping[tag], key)
        # else:
        #     mapping[tag] = key


def walk(node, path, mapping):
    try:
        add_mapping(node.tag, node.attrib, path, mapping)
    except ConflictMappingError as error:
        print('>', error.message)

    for child in node:
        walk(child, path, mapping)


def process_mapping(path, mapping):
    try:
        tree = ElementTree.parse(path)
        walk(tree.getroot(), path, mapping)
    except ElementTree.ParseError as error:
        global parsing_error_count
        parsing_error_count += 1
        print('Parsing error in {}: {}'.format(path, error))


def get_mapping(path):
    print('[Mapping] Processing mapping in {} ... '.format(path), flush=True)
    mapping = {}

    all_xmls = utility.get_all_xml_files(path)
    total = len(all_xmls)
    print('[Mapping] {} xmls to process'.format(total), flush=True)

    step = current = 0.05
    for index, file in enumerate(all_xmls):
        process_mapping(file, mapping)
        if (index + 1) / total >= current:
            print('[Mapping] {:.0f}% completed'.format(current * 100), flush=True)
            current += step

    global parsing_error_count
    print('[Mapping] Done! ({} parsing error)'.format(parsing_error_count), flush=True)
    return mapping


def dump(mapping, mapping_path):
    print('[Mapping] Dumping mapping to {} ...'.format(mapping_path), flush=True)
    string = json.dumps(mapping)
    with open(mapping_path, 'w') as file:
        file.write(string)
    print('[Mapping] Done!', flush=True)
