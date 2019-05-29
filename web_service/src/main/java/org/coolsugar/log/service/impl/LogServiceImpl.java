package org.coolsugar.log.service.impl;

import org.coolsugar.log.po.User;
import org.coolsugar.log.service.LogService;
import org.coolsugar.log.service.UserService;
import org.coolsugar.log.vo.OutUser;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.context.support.ClassPathXmlApplicationContext;
import org.springframework.stereotype.Service;

import javax.servlet.http.HttpServletRequest;


@Service
public class LogServiceImpl implements LogService {

    @Autowired
    private UserService userService;


    @Override
    public void sayHi() {

        ClassPathXmlApplicationContext context = new ClassPathXmlApplicationContext(new String[]{"classpath:applicationContext-MVC.xml"});

//        UserService userService = (UserService) context.getBean("userService");
        userService.sayHi();

        System.out.println("sending data");
    }


    @Override
    public void saveUser() {
        ClassPathXmlApplicationContext context = new ClassPathXmlApplicationContext(new String[]{"classpath:applicationContext-MVC.xml"});

        UserService userService = (UserService) context.getBean("userService");

        OutUser outUser = new OutUser(); //外部传进来的数据


        //测试更新功能
        outUser.setId(1);
        outUser.setCreator(0);
        outUser.setName("我的天");

        User user = new User();
        user.setCreator(outUser.getCreator());
        user.setName(outUser.getName());

        userService.saveUser(user);
    }


}
