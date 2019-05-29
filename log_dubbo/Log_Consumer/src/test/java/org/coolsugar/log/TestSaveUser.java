package org.coolsugar.log;

import org.coolsugar.log.po.User;
import org.coolsugar.log.vo.OutUser;
import org.coolsugar.log.service.UserService;
import org.springframework.context.support.ClassPathXmlApplicationContext;

import java.io.IOException;

/**
 * 保存功能，包括增加和修改
 * 2018/7/8-9
 */
public class TestSaveUser {
    public static void main(String[] args) {
        ClassPathXmlApplicationContext context = new ClassPathXmlApplicationContext(new String[]{"classpath:consumer.xml"});

        context.start();

        UserService userService = (UserService) context.getBean("userService");

        OutUser outUser = new OutUser(); //外部传进来的数据

        //测试插入功能
        outUser.setCreator(0);
        outUser.setName("hahaah");

        User user = new User();
        user.setCreator(outUser.getCreator());
        user.setName(outUser.getName());
        //其他字段属于默认字段，字段在sql语句中插入

        //测试更新功能
        /*outUser.setId(1);
        outUser.setCreator(0);
        outUser.setName("我的天");

        User user = new User();
        user.setId(outUser.getId());
        user.setCreator(outUser.getCreator());
        user.setName(outUser.getName());
*/
        userService.saveUser(user);

        System.out.println("sending data");

        try {
            System.in.read(); //按任意键结束
        } catch (IOException e) {
            e.printStackTrace();
        }
    }
}
