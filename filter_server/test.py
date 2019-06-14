import requests
import json

kv={"key":"1111","content":"32323231","type":"set"}
r = requests.request("GET", "http://127.0.0.1:7701/index", params=kv)
data = json.loads(r.text)
print(data)
