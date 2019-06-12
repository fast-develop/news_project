import tornado.ioloop
import tornado.web
import json
from elasticsearch import Elasticsearch
import time


class MainHandler(tornado.web.RequestHandler):
    def get_list_from_es(self, keyword):
        try:
            es = Elasticsearch(["127.0.0.1:8010"])
            json_body = {
                    "query":
                        {
                         "match":
                            {
                             "keyword":keyword
                            }
                        }
                    }
            rsp_body = es.search(index="doc", doc_type="keyword", body=json_body)
            print(rsp_body)

            result = {}
            result["keyword"] = keyword
            result["doclist"] = []
            for k in rsp_body["hits"]["hits"]:
                tmp = {}
                tmp["docid"] = k["_source"]["docid"] 
                tmp["score"] = k["_score"]
                tmp["from"] = "keyword:"+ keyword
                result["doclist"].append(tmp)

            return result
        except Exception as e:
            print(e)
            return {}

    def get(self):
        args=self.request.arguments
        print(args)
        
        if "keyword" in args:
            result = {}
            result["err_id"] = 0
            result["data"] = []
            keywords = args["keyword"]
            for keyword in keywords:
                try:
                    key = keyword.decode('utf-8')
                    re = self.get_list_from_es(key)
                    result["data"].append(re)
                except Exception as e:
                    print(e)
            self.write(json.dumps(result, ensure_ascii=False))
            return

        else:
            print("no keyword")
            result = {}
            result["err_id"] = -1
            result["err_msg"] = "no keyword"
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
    application.listen(7501)
    # epoll + socket
    tornado.ioloop.IOLoop.instance().start()

