package org.coolsugar.web.controller;

import org.coolsugar.web.service.NewsService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;

import javax.annotation.Resource;
import javax.imageio.ImageIO;
import javax.servlet.ServletOutputStream;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.awt.image.BufferedImage;
import java.io.IOException;
import java.io.PrintWriter;
import java.util.List;

@Controller
@RequestMapping("/news")
public class NewsController {

    @Autowired
    private NewsService service;

    @RequestMapping("/index")
    public String getNewsIndex() {
        return "news_index";
    }

    @RequestMapping("/list")
    public void getNewsList(HttpServletResponse response) throws IOException {

        List<String> ret = service.getNewsList();

        response.setCharacterEncoding("UTF-8");
        response.setContentType("application/json; charset=utf-8");

        PrintWriter pw = response.getWriter();
        pw.write(ret.toString());
    }

    @RequestMapping("/detail/{id}")
    protected void getNewsDetail(HttpServletResponse response, @PathVariable String id) throws IOException {

        String ret = service.getNewsDetail(id.substring(12));

        response.setCharacterEncoding("UTF-8");
        response.setContentType("application/json; charset=utf-8");

        PrintWriter pw = response.getWriter();
        pw.write(ret);
    }

    @RequestMapping("/thumb")
    protected void getThumb(HttpServletRequest request, HttpServletResponse response) throws IOException {

        String thumb = request.getParameter("thumburl");

        BufferedImage buffImg = service.getThumb(thumb);

        response.setContentType("image/png");

        ServletOutputStream sos = response.getOutputStream();
        ImageIO.write(buffImg,"png", sos);
        sos.close();
    }
}
