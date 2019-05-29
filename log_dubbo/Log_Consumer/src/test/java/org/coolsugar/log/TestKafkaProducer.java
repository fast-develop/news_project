package org.coolsugar.log;

import org.coolsugar.log.service.LogService;
import org.springframework.context.support.ClassPathXmlApplicationContext;

public class TestKafkaProducer {

    public static void main(String[] args) {
        ClassPathXmlApplicationContext context = new ClassPathXmlApplicationContext(new String[]{"classpath:consumer.xml"});
        context.start();

        LogService logService = (LogService) context.getBean("logService");
        logService.sendMsg();

    }
}
