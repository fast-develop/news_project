import json
import time
from elasticsearch import Elasticsearch

class RandomHandler:
    def get(self):
        print("random\n")
        es = Elasticsearch(["119.3.252.18:8010"])

        cur_time = str(int(time.time()))
        param_body={
           "size": 10,
           "query": {
              "function_score": {
                 "functions": [
                    {
                       "random_score": {
                          "seed": cur_time
                       }
                    }
                 ]
              }
           }
        }
        rsp_body = es.search(index="doc", doc_type="keyword", body=param_body)

        result = []
        for k in rsp_body["hits"]["hits"]:
            tmp = {}
            tmp["docid"] = k["_source"]["docid"] 
            tmp["score"] = 0
            tmp["from"] = "random"
            result.append(tmp)
    
        #print("random list:",result)

        return result 
