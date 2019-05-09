import sys
from pymongo import MongoClient

class MyMongoClient:
    def __init__(self,ip,p,db,table):
        self.conn = MongoClient(host=ip, port=p)
	self.db = self.conn[db]
        self.collection = self.db[table]

    def get_random(self):
        for item in self.collection.aggregate([ {'$sample': {'size':10}}]):
            print item['_id']


if __name__ == "__main__":
    mongo_client = MyMongoClient('127.0.0.1',27017,'meta','toutiao')  
    mongo_client.get_random()
