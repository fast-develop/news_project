package org.coolsugar.log;

import org.coolsugar.log.service.UserService;
import org.springframework.context.support.ClassPathXmlApplicationContext;

public class TestSayHi {

    public static void main(String[] args) {
        ClassPathXmlApplicationContext context = new ClassPathXmlApplicationContext(new String[]{"classpath:consumer.xml"});

        context.start();

        UserService userService = (UserService) context.getBean("userService");


        userService.sayHi();

    }
}
