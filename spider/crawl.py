#!coding=utf-8
import requests
import os
import re
import io
import json
import math
import random
import time
import HTMLParser
from requests.packages.urllib3.exceptions import InsecureRequestWarning
import pandas as pd
requests.packages.urllib3.disable_warnings(InsecureRequestWarning)  ###禁止提醒SSL警告
import hashlib
import execjs
import datetime
from bs4 import BeautifulSoup
from pymongo import MongoClient


class toutiao(object):

    def __init__(self):
        self.s = requests.session()
        headers = {'Accept': '*/*',
                   'Accept-Language': 'zh-CN',
                   'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; WOW64; Trident/7.0; .NET4.0C; .NET4.0E; .NET CLR 2.0.50727; .NET CLR 3.0.30729; .NET CLR 3.5.30729; InfoPath.3; rv:11.0) like Gecko',
                   'Connection': 'Keep-Alive',

                   }
        self.s.headers.update(headers)

        self.conn = MongoClient('127.0.0.1', 27017)
        self.db = self.conn.pages
        self.my_set = self.db.toutiao



    def closes(self):
        self.s.close()

    def write_line(self, fd, docid, url):
        line = str(docid) + "\t" + url + "\n"
        print line
        fd.write(line)

    def getdata(self, url):  #获取数据
        req = self.s.get(url=url, verify=False)
        soup = BeautifulSoup(req.text,"html.parser")
        #data = soup.find_all("script")[6].string[16:-1]
        data = soup.find_all("script")[6].text
	reldatematch = re.search(r'articleInfo: {(.|\s)+?}', data).group()
        title = re.search(r'title: \'(.*?)\'', reldatematch).group(1)
        content = re.search(r'content: \'(.*?)\'', reldatematch).group(1)
	content = HTMLParser.HTMLParser().unescape(content)
        #print content
        docid = re.search(r'itemId: \'(.*?)\'', reldatematch).group(1)
        print docid
	now=time.time()
	crawl_time = time.strftime("%Y-%m-%d %H:%M:%S",time.localtime(now))  ##抓取时间
	self.my_set.insert({"docid":docid, "title":title, "content":content, "time":crawl_time})
        
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
        f = io.open("/data/app/spider/a.js", 'r', encoding='UTF-8')
        line = f.readline()
        htmlstr = ''
        while line:
            htmlstr = htmlstr + line
            line = f.readline()
        try:
            ctx = execjs.compile(htmlstr)
        except Exception as e:
            print e
        f.close()
        return ctx.call('get_as_cp_signature')

if __name__=='__main__':
    tt = toutiao()
    path = "/data/app/spider/data/"
    cur_file = (datetime.datetime.now()-datetime.timedelta(hours=1)).strftime("%Y-%m-%d-%H")
    file_path = path + cur_file
    print file_path
    try:
        fd = open(file_path, "r")
    except Exception as e:
        print e
        os._exit(1)

    for line in fd:
        try:
            items = line.split("\t")
            docid = items[0]
            url = items[1]
    	    tt.getdata(url)
            time.sleep(1)
        except Exception as e:
            print e



