package org.coolsugar.log;

import org.coolsugar.log.po.User;
import org.coolsugar.log.service.UserService;
import org.coolsugar.log.vo.OutUser;
import org.springframework.context.ApplicationContext;
import org.springframework.context.support.ClassPathXmlApplicationContext;

public class TestLog {

    private UserService userService;

    public static void main(String[] args) {

        TestLog tl = new TestLog();
        tl.testInsert();
    }

    private void testInsert( ) {

        ApplicationContext context = new ClassPathXmlApplicationContext("dubbo-consumer.xml");

        userService = (UserService) context.getBean("userService");

        //测试插入功能
        OutUser outUser = new OutUser(); //外部传进来的数据
        outUser.setCreator(0);
        outUser.setName("john");

        User user = new User();
        user.setCreator(outUser.getCreator());
        user.setName(outUser.getName());
        userService.saveUser(user);

        System.out.println("sending data");
    }
}
