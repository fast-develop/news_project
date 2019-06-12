import requests
import json

kv={"keyword":["5G","中国"]}
r = requests.request("GET", "http://127.0.0.1:7501/index", params=kv)
data = json.loads(r.text)
if data["err_id"] != 0:
    print("error:", data["err_msg"])
else:
    print(data)
    doc_list = []
    for item in data["data"]:
        print(item["keyword"])
        for i in item["doclist"]:
            doc_list.append(i)
