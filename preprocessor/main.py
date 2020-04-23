import argparse
import os

parser = argparse.ArgumentParser(description='ElasticJury Data Preprocessor')
parser.add_argument('--path', type=str, default='demo', help='Relative path for data to process')


def process(path):
    pass


def run(path):
    for home, dirs, files in os.walk(path):
        for file in files:
            process(file)


if __name__ == '__main__':
    args = parser.parse_args()
    run(args.path)
