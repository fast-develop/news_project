from bs4 import BeautifulSoup
import re
import json


html = '''<script type="text/javascript">

        var TralbumData = {

            current: {
                "upc":null,"title":"BLACK&WHITE MEDICINE","purchase_title":null,
                "download_desc_id":null,"minimum_price":0.0,"set_price":7.0,"mod_date":"17 Jun 2018 11:47:50 GMT"
            },
            album_is_preorder: null,
            album_release_date: "17 Jun 2018 00:00:00 GMT",
        };
        </script>'''



#pattern = re.compile(r'album_release_date: \"(.*?)\"', re.MULTILINE)
#pattern = re.compile(r'current: {(.|\s)+?}', re.MULTILINE)
pattern = re.compile(r'var TralbumData = {(.|\s)+?};', re.MULTILINE)
soup = BeautifulSoup(html, 'html.parser')
script = soup.script.text
release_date = re.search(pattern, script).group(0)[17:-1]
print release_date
json_d = json.loads(release_date)
print (release_date)
