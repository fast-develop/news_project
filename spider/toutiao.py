#!coding=utf-8
import requests
import re
import io
import json
import math
import random
import time
from requests.packages.urllib3.exceptions import InsecureRequestWarning
import pandas as pd
requests.packages.urllib3.disable_warnings(InsecureRequestWarning)  ###禁止提醒SSL警告
import hashlib
import execjs
import datetime
from pymongo import MongoClient


class toutiao(object):

    def __init__(self,path,url):
        self.path = path  # CSV保存地址
        self.url=url
        self.s = requests.session()
        headers = {'Accept': '*/*',
                   'Accept-Language': 'zh-CN',
                   'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; WOW64; Trident/7.0; .NET4.0C; .NET4.0E; .NET CLR 2.0.50727; .NET CLR 3.0.30729; .NET CLR 3.5.30729; InfoPath.3; rv:11.0) like Gecko',
                   'Connection': 'Keep-Alive',

                   }
        self.s.headers.update(headers)
        self.channel=re.search('ch/(.*?)/',url).group(1)

        self.conn = MongoClient('127.0.0.1', 27017)
        self.db = self.conn.meta
        self.my_set = self.db.toutiao



    def closes(self):
        self.s.close()

    def write_line(self, fd, docid, url):
        line = str(docid) + "\t" + url + "\n"
        print line
        fd.write(line)

    def getdata(self):  #获取数据
        path = "/data/app/spider/data/"
        thumb_path = "/data/app/spider/thumb_data/"
        cur_file = (datetime.datetime.now()-datetime.timedelta()).strftime("%Y-%m-%d-%H")
        file_path = path + cur_file
        thumb_file_path = thumb_path + cur_file
        print file_path
        fd = open(file_path, "a+")
        thumb_fd = open(thumb_file_path, "a+")
        req = self.s.get(url=self.url, verify=False)
        #print (self.s.headers)
        #print(req.text)
        headers = {'referer': self.url}
        max_behot_time='0'
        signature='.1.hXgAApDNVcKHe5jmqy.9f4U'
        eas = 'A1E56B6786B47FE'
        ecp = '5B7674A7FF2E9E1'
        self.s.headers.update(headers)
	'''
        title = []
        source = []
        source_url = []
        comments_count = []
        tag = []
        chinese_tag = []
        label = []
        abstract = []
        behot_time = []
        nowtime = []
        duration = []
        '''	
        for i in range(0,30):  ##获取页数
            try:
		    Honey = json.loads(self.get_js())
		    eas = Honey['as']
		    ecp = Honey['cp']
		    signature = Honey['_signature']
		    url='https://www.toutiao.com/api/pc/feed/?category={}&utm_source=toutiao&widen=1&max_behot_time={}&max_behot_time_tmp={}&tadrequire=true&as={}&cp={}&_signature={}'.format(self.channel,max_behot_time,max_behot_time,eas,ecp,signature)
		    req=self.s.get(url=url, verify=False)
		    time.sleep(random.random() * 2+2)
		    #print(req.text)
		    print(url)
		    j=json.loads(req.text)

		    for k in range(0, 10):
			try:
				now=time.time()
				if j['data'][k]['tag'] != 'ad':
				    title = j['data'][k]['title'] ##标题
				    docid = j['data'][k]['item_id']
				    author = j['data'][k]['source'] ##作者
				    url = 'https://www.toutiao.com'+j['data'][k]['source_url']  ##文章链接
				    tag = j['data'][k]['tag']  ###频道名
				    abstract = j['data'][k]['abstract'] ###文章摘要
				    image_url = 'https:'+j['data'][k]['image_url']
				    crawl_time = time.strftime("%Y-%m-%d %H:%M:%S",time.localtime(now))  ##抓取时间
				    print docid, url,image_url
				    self.write_line(fd, docid, url)             
                                    self.write_line(thumb_fd, docid, image_url)
				    self.my_set.insert({"docid":docid, "url":url, "tag":tag, "abstract":abstract,"title":title,"author":author,"time":crawl_time, "image_url":image_url})
				    time.sleep(1)
			except Exception as e:
			    print e

		    time.sleep(2)

		    #max_behot_time=str(j['next']['max_behot_time'])
		    '''
		    print('------------'+str(j['next']['max_behot_time']))
		    print(title)
		    print(source)
		    print(source_url)
		    print(comments_count)
		    print(tag)
		    print(chinese_tag)
		    print(label)
		    print(abstract)
		    print(behot_time)
		    print(nowtime)
		    print(duration)
		    '''
		    #print(source_url)
            except Exception as e:
                print e

        fd.close()

    def getHoney(self,t):  #####根据JS脚本破解as ,cp
        #t = int(time.time())  #获取当前时间
        #t=1534389637
        #print(t)
        e =str('%X' % t)  ##格式化时间
        #print(e)
        m1 = hashlib.md5()  ##MD5加密
        m1.update(str(t).encode(encoding='utf-8'))  ##转化格式
        i = str(m1.hexdigest()).upper() ####转化大写
        #print(i)
        n=i[0:5]    ##获取前5位
        a=i[-5:]    ##获取后5位
        s=''
        r=''
        for x in range(0,5):
            s+=n[x]+e[x]
            r+=e[x+3]+a[x]
        eas='A1'+ s+ e[-3:]
        ecp=e[0:3]+r+'E1'
        #print(eas)
        #print(ecp)
        return eas,ecp

    def get_js(self):  ###二牛破解as ,cp,  _signature  参数的代码，然而具体关系不确定，不能连续爬取
        # f = open("D:/WorkSpace/MyWorkSpace/jsdemo/js/des_rsa.js",'r',encoding='UTF-8')
        print "eeee"
        f = io.open("/data/app/spider/a.js", 'r', encoding='UTF-8')
        print "ffff"
        line = f.readline()
        htmlstr = ''
        while line:
            htmlstr = htmlstr + line
            line = f.readline()
        try:
            ctx = execjs.compile(htmlstr)
        except Exception as e:
            print e
        print "xxx"
        f.close()
        return ctx.call('get_as_cp_signature')

if __name__=='__main__':
    url='https://www.toutiao.com/ch/news_tech/'
    tt = toutiao("./b", url)
    while True:
        tt.getdata()


