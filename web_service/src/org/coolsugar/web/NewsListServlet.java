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

import java.util.List;


@WebServlet(name = "NewsListServlet")
public class NewsListServlet extends HttpServlet {
    protected void doPost(HttpServletRequest request, HttpServletResponse response) throws ServletException, IOException {

        NewsService service = BaseFactory.getFactory().getInstance(NewsService.class);
        List<String> ret = service.getNewsList();

        response.setCharacterEncoding("UTF-8");
        response.setContentType("application/json; charset=utf-8");

        PrintWriter pw = response.getWriter();
        pw.write(ret.toString());
    }

    protected void doGet(HttpServletRequest request, HttpServletResponse response) throws ServletException, IOException {
        this.doPost(request, response);
    }
}
