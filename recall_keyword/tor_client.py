import tornado.httpclient

http_client = tornado.httpclient.HTTPClient()
try:
    response = http_client.fetch("http://127.0.0.1:7051/index?keyword=xuke")
    print(response.body)
except Exception as e:
    print(e)
