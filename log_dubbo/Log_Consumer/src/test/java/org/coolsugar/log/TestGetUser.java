package org.coolsugar.log;

import org.coolsugar.log.po.User;
import org.coolsugar.log.vo.OutUser;
import org.coolsugar.log.service.UserService;
import org.springframework.context.support.ClassPathXmlApplicationContext;

import java.io.IOException;

/**
 * 根据id查找用户的功能
 * 2018/7/8
 */
public class TestGetUser {
    public static void main(String[] args) {
        ClassPathXmlApplicationContext context = new ClassPathXmlApplicationContext(new String[] { "classpath:consumer.xml" });

        context.start();

        UserService userService = (UserService) context.getBean("userService");
        User user = userService.getUserById(1);  //PO
        //System.out.println(user.toString());

        OutUser outUser  = new OutUser();  //VO
        outUser.setCreator(user.getCreator());
        outUser.setId(user.getId());
        outUser.setName(user.getName());
        System.out.println(outUser.toString()); //返回给用户的数据

        try {
            System.in.read(); //按任意键结束
        } catch (IOException e) {
            e.printStackTrace();
        }
    }
}
