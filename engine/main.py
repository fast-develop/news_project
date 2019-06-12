import tornado.ioloop
import tornado.web
import json
import time
import redis
import requests
from recall_keyword import KeywordHandler
from recall_random import RandomHandler

recall_map = {
    "random":RandomHandler,
    "keyword":KeywordHandler,
}

class MainHandler(tornado.web.RequestHandler):
    def filter_doc_list(self, userid, doc_list):
        new_doc_list = []
        for doc in doc_list:
            try:
                docid = doc['docid']
                kv={"key":userid,"content":docid,"type":"get"}
                r = requests.request("GET", "http://127.0.0.1:7701/index", params=kv)
                data = json.loads(r.text)
                if data['exist'] == False:
                    new_doc_list.append(doc)
                else:
                    print("filter", docid)
            except Exception as e:
                print(e)

        return new_doc_list
                    
    def get_user_click_list(self,userid):
        try:
            r = redis.Redis(host='127.0.0.1',port=6379, db=1)
            re = r.lrange(userid,0,10)[0].decode()
            re = re.strip('[]')
            re = re.split(', ')
            result = []
            for k in re:
                try:
                    items = k.split(':')
                    docid = items[0].strip()
                    time = int(items[1].strip())
                    item = [docid,time]
                    result.append(item)
                except Exception as e:
                    print(e)
            return result
        except Exception as e:
            print(e)
            return []

    def get(self):
        args=self.request.arguments
        print(args)
        cur_time = int(time.time())
        
        if "userid" in args:
            result = {}
            result["err_id"] = 0
            result["data"] = []
            print(args["userid"])
            userid = args["userid"][0].decode('utf-8')
            item_list = self.get_user_click_list(userid)
                   
            doc_list = []
            
            #get recall doclist from recall_map handler
            for k,v in recall_map.items():
                doc_list.extend(v.get(item_list))

            if len(doc_list) != 0:
                #sort
                doc_list.sort(key=lambda k: (k.get('score',0)), reverse=True) 

            #filter
            final_doc_list = self.filter_doc_list(userid, doc_list)

            result["data"] = final_doc_list[0:5]
            self.write(json.dumps(result, ensure_ascii=False))
            return

        else:
            print("userid")
            result = {}
            result["err_id"] = -1
            result["err_msg"] = "no userid"
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
    keywordhandler = KeywordHandler()
    recall_map["keyword"] = keywordhandler

    application.listen(7602)
    # epoll + socket
    tornado.ioloop.IOLoop.instance().start()

    
