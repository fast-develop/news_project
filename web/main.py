import tornado.ioloop
import tornado.web
import json
from pymongo import MongoClient
from cStringIO import StringIO
from PIL import Image


class MyMongoClient:
    def __init__(self,ip,p,db,table):
        self.conn = MongoClient(host=ip, port=p)
        self.db = self.conn[db]
        self.collection = self.db[table]

    def get_random(self):
        return self.collection.aggregate([{'$sample': {'size':10}}])

    def find_thumb(self,url):
        return self.collection.find({'thumb_url':url})


class MainHandler(tornado.web.RequestHandler):
    def get(self):
        mongo_client = MyMongoClient('127.0.0.1',27017,'meta','toutiao')
        list = mongo_client.get_random()

        args=self.request.arguments
	print args
	self.set_header('Content-Type', 'application/json; charset=UTF-8')

        result = []
        for item in list:
            i_dict = {}
            i_dict['title'] = item['title']
            i_dict['_id'] = item['docid']
            i_dict['brief'] = 'brief'
            i_dict['category'] = 'category'
            i_dict['link'] = item['url']
            i_dict['thumb'] = item['image_url']
            i_dict['publisher'] = item['author']
            i_dict['pubData'] = item['time']
            result.append(i_dict)

	    print i_dict
            
	self.write(json.dumps(result, ensure_ascii=False))

    def post(self, *args, **kwargs):
        self.write('post')

    def delete(self):
        self.write("delete")

class ThumbHandler(tornado.web.RequestHandler):
    def get(self):
        args=self.request.arguments
	print "thumb"
	thumb_url = args['thumburl'][0] + '\n'
        print thumb_url
	self.set_header('Content-Type', 'image/png')
	mongo_client = MyMongoClient('127.0.0.1',27017,'meta','thumb')
	data = mongo_client.find_thumb(thumb_url)
        for item in list(data):
            image = item['data']
	    self.write(image)

    def post(self, *args, **kwargs):
        self.write('post')

    def delete(self):
        self.write("delete")

class ArticleHandler(tornado.web.RequestHandler):
    def get(self,id):
        args=self.request.arguments
	print "article",id
	print args
	self.set_header('Content-Type', 'application/json; charset=UTF-8')
	result={"id":"111","text":"xxxfewfewfewfewfsfdsfsdffafeowfjewofnifewfewf"}	
	self.write(json.dumps(result, ensure_ascii=False))

    def post(self, *args, **kwargs):
        self.write('post')

    def delete(self):
        self.write("delete")


settings = {
    'template_path': 'template',
    'static_path': 'static',
}


application = tornado.web.Application([
    (r"/index", MainHandler),
    (r"/thumb", ThumbHandler),
    (r"/article/(?P<id>\d*)", ArticleHandler),
], **settings)

if __name__ == "__main__":
    application.listen(8080)
    # epoll + socket
    tornado.ioloop.IOLoop.instance().start()

    mongo_client = MyMongoClient('127.0.0.1',27017,'meta','toutiao')
    mongo_client.get_random()


