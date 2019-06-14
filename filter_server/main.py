import tornado.ioloop
import tornado.web
import json
from elasticsearch import Elasticsearch
import time
from bloomfilter import BloomFilter

g_filter = BloomFilter()

class MainHandler(tornado.web.RequestHandler):
    def get(self):
        args=self.request.arguments
        print(args)
        
        if "key" in args and "content" in args and "type" in args:
            tp = args['type'][0].decode('utf-8')
            result = {}
            result["err_id"] = 0

            key = args['key'][0].decode('utf-8')
            content = args['content'][0].decode('utf-8')

            if tp == "get":
                if g_filter.is_contains(key, content):
                    print(content + ' is existed')
                    result['exist'] = True
                else:
                    print(content + ' is not existed')
                    result['exist'] = False

                self.write(json.dumps(result, ensure_ascii=False))
                return
            else:
                result['set'] = 'sucess'
                g_filter.insert(key, content)

            self.write(json.dumps(result, ensure_ascii=False))
            return

        else:
            print("no key or content")
            result = {}
            result["err_id"] = -1
            result["err_msg"] = "no key or content"
            self.write(json.dumps(result, ensure_ascii=False))
            return

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
], **settings)

if __name__ == "__main__":
    application.listen(7701)
    # epoll + socket
    tornado.ioloop.IOLoop.instance().start()

