package org.coolsugar.web;

import org.coolsugar.log.po.User;
import org.coolsugar.log.service.UserService;
import org.coolsugar.log.vo.OutUser;
import org.springframework.context.ApplicationContext;
import org.springframework.context.support.ClassPathXmlApplicationContext;

public class TestLog {

    private UserService userService;

    public static void main(String[] args) {

        TestLog tl = new TestLog();
        tl.sayHi();
    }

    private void sayHi( ) {

        ApplicationContext context = new ClassPathXmlApplicationContext("dubbo-consumer.xml");

        userService = (UserService) context.getBean("userService");

        userService.sayHi();

        System.out.println("sending data");
    }
}
