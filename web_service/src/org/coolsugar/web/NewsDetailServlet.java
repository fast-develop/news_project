package org.coolsugar.web;

import org.coolsugar.service.NewsService;
import org.coolsugar.util.BaseFactory;

import javax.servlet.ServletException;
import javax.servlet.annotation.WebServlet;
import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.io.IOException;
import java.io.PrintWriter;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

@WebServlet(name = "NewsDetailServlet")
public class NewsDetailServlet extends HttpServlet {
    protected void doPost(HttpServletRequest request, HttpServletResponse response) throws ServletException, IOException {
        String url = request.getRequestURI();
        Pattern pattern = Pattern.compile("/.+/.+/\\d{13}(.*)");// 匹配的模式
        Matcher m = pattern.matcher(url);
        String id = m.find() ? m.group(1) : "";

        NewsService service = BaseFactory.getFactory().getInstance(NewsService.class);
        String ret = service.getNewsDetail(id);

        response.setCharacterEncoding("UTF-8");
        response.setContentType("application/json; charset=utf-8");

        PrintWriter pw = response.getWriter();
        pw.write(ret);
    }

    protected void doGet(HttpServletRequest request, HttpServletResponse response) throws ServletException, IOException {
        doPost(request, response);
    }

}
