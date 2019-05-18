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
  </head>
  <body>
    <div id="content">
      <div id="search_div">
        <form method="post" action="">
          <span class="text-input">关键词: <input type="text" name="keywords"/> </span>
          <input type="submit" value="筛选">
        </form>
      </div>
      <div id="news_list">
        <c:forEach items="${requestScope.news_list}" var="news">
          <div id="news">
            <img src="thumb?thumburl=https://p1.pstatp.com/list/190x124/pgc-image/324aed32d85f40c4a4d2b556c0ff808d"/> <a href="">news</a>
          </div>
          
        </c:forEach>
      </div>
    </div>
  </body>
</html>
