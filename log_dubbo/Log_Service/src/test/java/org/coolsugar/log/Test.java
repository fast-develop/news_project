package org.coolsugar.log;

import org.springframework.context.support.ClassPathXmlApplicationContext;

import java.io.IOException;

public class Test {
    public static void main(String[] args) throws Exception {
        ClassPathXmlApplicationContext context = new ClassPathXmlApplicationContext("classpath:provider.xml");
        context.start();

        System.out.println("dubbo provider starting...");
        try {
            System.in.read(); //按任意键退出
        } catch (IOException e) {
            e.printStackTrace();
        }
    }
}
