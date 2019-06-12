import json
import time
import requests
from elasticsearch import Elasticsearch

class KeywordHandler():
    def get(self, click_list):
        print("key\n")
        doc_list = []
        cur_time = int(time.time())
        for item in click_list:
            try:
                if cur_time - item[1] < 300: #over five minutes 
                    doc_list.append(item[0])         
                else:
                    print("over fine minutes")

            except Exception as e:
                print(e)

        doc_list = list(set(doc_list))#去重
        #print(doc_list)

        if len(doc_list) == 0:
            print("no recent read")
            return []
            
        re = self.get_doc_from_es(doc_list)                                                                                                        
        if len(re) == 0:
            print("get nothing")
            return []

        keywords = []
        for item in re: 
            for i in item['keyword']:
                #keywords = keywords.extend(item['keyword'])
                keywords.append(i)
        #print(keywords)
        #keywords = list(set(keywords))#去重

        kv = {"keyword":keywords}
        print(kv)
        rsp = requests.request("GET", "http://127.0.0.1:7501/index", params=kv)
        data = json.loads(rsp.text)
        if data["err_id"] != 0:
            print("error:", data["err_msg"])
            return []
        else:
            doc_list = []
            for item in data["data"]:
                for i in item["doclist"]:
                    doc_list.append(i)    
            #print("keyword list:",doc_list)
            return doc_list

    def get_doc_from_es(self, doc_id_list):
        es = Elasticsearch(["127.0.0.1:8010"])
        json_body = {
                "query":
                    {
                     "constant_score":
                        {
                         "filter":
                            {
                             "terms":
                                {
                                 "docid":doc_id_list
                                }
                            }
                        }
                    }
                }
        rsp_body = es.search(index="doc", doc_type="keyword", body=json_body)
        #print(rsp_body)

        result = []
        for k in rsp_body["hits"]["hits"]:
            tmp = {}
            tmp["docid"] = k["_source"]["docid"] 
            tmp["keyword"] = k["_source"]["keyword"]
            result.append(tmp)

        return result


