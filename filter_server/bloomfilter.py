# coding: utf-8

import mmh3
import redis


BIT_SIZE = 5000
SEEDS = [50, 51, 52, 53, 54, 55, 56]


def get_redis(host='localhost', port=6379, db=10):
    return redis.Redis(host=host, port=port, db=db)


class BloomFilter(object):

    def __init__(self):
        self.db = get_redis()

    def cal_offsets(self, content):
        return [mmh3.hash(content, seed) % BIT_SIZE for seed in SEEDS]

    def is_contains(self, key, content):
        if not content:
            return False
        locs = self.cal_offsets(content)

        return all(True if self.db.getbit(key, loc) else False for loc in locs)

    def insert(self, key, content):
        locs = self.cal_offsets(content)
        for loc in locs:
            self.db.setbit(key, loc, 1)


'''
if __name__ == '__main__':
    bloom_filter = BloomFilter()

    test_url = 'https://douban.com1'

    print('before')
    if bloom_filter.is_contains('1111',test_url):
        print(test_url + ' is existed')
    else:
        print(test_url + ' is not existed')

    bloom_filter.insert('1111',test_url)

    print('after')
    if bloom_filter.is_contains('1111',test_url):
        print(test_url + ' is existed')
    else:
        print(test_url + ' is not existed')
'''
