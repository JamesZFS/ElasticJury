# -*- coding: UTF-8 -*-

import argparse
import json

import mapping_collector
import processor

from utility import *

parser = argparse.ArgumentParser(description='ElasticJury Data Preprocessor')
parser.add_argument('--path', type=str, default='../dataset/xml_1', help='Relative path for data to process')
parser.add_argument('--mapping', type=str, default='mapping.json', help='Tag to nameCN mapping file path')
parser.add_argument('--password', type=str, default='', help='Password to login your local mysql account')
parser.add_argument('--no-db', dest='do_db', action='store_false', help='Whether run database processing')
parser.add_argument('--no-mapping', dest='do_mapping', action='store_false', help='Whether run mapping')
parser.add_argument('--clean-mapping', dest='clean_mapping', action='store_true', help='Clean the mapping file')
parser.set_defaults(do_db=True, do_mapping=True, clean_mapping=False)


def run(path, mapping_path, db_password, do_db, do_mapping, clean_mapping):
    # Check if xml path exists
    if not os.path.exists(path):
        log_exit('Main', 'XML data directory path {} does not exist'.format(path))

    # Clean mapping
    if os.path.exists(mapping_path) and clean_mapping:
        os.remove(mapping_path)

    # Create mapping if not exists
    if (not os.path.exists(mapping_path)) and do_mapping:
        mapping = mapping_collector.get_mapping(path)
        mapping_collector.dump(mapping, mapping_path)

    # Read mapping
    if not os.path.exists(mapping_path):
        log_exit('Main', 'Mapping file {} does not exist'.format(mapping_path))
    with open(mapping_path, 'r') as file:
        mapping = json.load(file)

    # Processing and dump into database
    if do_db:
        processor.process(mapping, path, db_password)


if __name__ == '__main__':
    args = parser.parse_args()
    run(args.path, args.mapping, args.password, args.do_db, args.do_mapping, args.clean_mapping)
