import requests
import json

kv={"userid":"xuke"}
r = requests.request("GET", "http://127.0.0.1:7602/index", params=kv)
data = json.loads(r.text)
if data["err_id"] != 0:
    print("error:", data["err_msg"])
else:
    for i in data['data']:
        print(i)
