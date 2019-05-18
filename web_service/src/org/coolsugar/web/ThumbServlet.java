package org.coolsugar.web;

import org.coolsugar.service.NewsService;
import org.coolsugar.util.BaseFactory;

import javax.imageio.ImageIO;
import javax.servlet.ServletException;
import javax.servlet.ServletOutputStream;
import javax.servlet.annotation.WebServlet;
import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.awt.image.BufferedImage;
import java.io.FileInputStream;
import java.io.IOException;
import java.io.PrintWriter;

@WebServlet(name = "ThumbServlet")
public class ThumbServlet extends HttpServlet {
    protected void doPost(HttpServletRequest request, HttpServletResponse response) throws ServletException, IOException {

        String thumb = request.getParameter("thumburl");

        NewsService service = BaseFactory.getFactory().getInstance(NewsService.class);
        BufferedImage buffImg = service.getThumb(thumb);

        response.setContentType("image/png");

        ServletOutputStream sos = response.getOutputStream();
        ImageIO.write(buffImg,"png", sos);
        sos.close();
    }

    protected void doGet(HttpServletRequest request, HttpServletResponse response) throws ServletException, IOException {
        doPost(request, response);
    }
}
