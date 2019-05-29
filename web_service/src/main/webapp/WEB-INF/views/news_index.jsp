<%--
  Created by IntelliJ IDEA.
  User: xuyi
  Date: 2019-05-10
  Time: 16:52
  To change this template use File | Settings | File Templates.
--%>
<%@ page contentType="text/html;charset=UTF-8" language="java" %>
<html>
  <head>
    <title>news</title>
    <style type="text/css">
      div{
        width: 960px;
        height: 120px;
        /*border:1px solid red ;*/
        /*line-height: 60px;*/
      }
      div img{
        vertical-align: middle;
      }
    </style>
    <script src="https://code.jquery.com/jquery-3.3.1.min.js"></script>
  </head>
  <body onload="getNewsList()">
    <div id="content">
      <div id="search_div">
        <form method="post" action="">
          <span class="text-input">关键词: <input type="text" name="keywords"/> </span>
          <input type="submit" value="筛选">
        </form>
      </div>
      <div id="news_list">
        <!--<img src="thumb?thumburl=https://p3.pstatp.com/list/190x124/pgc-image/RIPA7A33Omthul"/>-->
      </div>
      <script type="text/javascript">
        function getNewsList() {
            //alert("hahahahaha");
            $.getJSON("<%=request.getContextPath()%>/news/list", function(data){
              $.each(data, function (i, value) {


                /*
                var span = '<span id="img_"'+ value['_id'] +'> </span>';
                $("#news_list").append(span);
                */

                var news_title = value['title'];
                var news_link = value['link'];
                var news_img_src = '<%=request.getContextPath()%>/news/thumb?thumburl='+value['thumb'];
                var news_img = '<img alt="pic is missing..." src="'+news_img_src+'" title="'+news_title+'" />';

                $("#news_list").append(news_img);

                $("#news_list").append('<span> <a href="'+news_link+'">'+news_title+'</a> </span>');
                $("#news_list").append("<br/> <hr/>")


              })
            });
        }
      </script>
    </div>
  </body>
</html>
