package org.coolsugar.log.controller;

import org.coolsugar.log.service.LogService;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.RequestMapping;

import javax.annotation.Resource;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.io.IOException;
import java.io.PrintWriter;


@Controller
@RequestMapping("/log")
public class LogController {

    @Resource
    private LogService service;

    @RequestMapping("/hi")
    public void sayHi(HttpServletRequest request, HttpServletResponse response) throws IOException {

        System.out.println("sending data");

        writeHTML(response);

        service.sayHi();

    }

    @RequestMapping("/save_user")
    public void getUser(HttpServletRequest request, HttpServletResponse response) throws IOException {
        service.saveUser();
    }

    private void writeHTML(HttpServletResponse response) throws IOException{
        // 写到网页上
        PrintWriter pw = response.getWriter();
        pw.write("<!DOCTYPE HTML PUBLIC \\\"-//W3C//DTD HTML 4.01 Transitional//EN\\\">");
        pw.write("<HTML>");
        pw.write("<HEAD>this is a head</HEAD>");
        pw.write("</br>");
        pw.write("<BODY>this is a body</BODY>");
        pw.write("</HTML>");
    }
}
