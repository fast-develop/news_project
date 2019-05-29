package org.coolsugar.log.impl;

import org.coolsugar.log.kafka.LogProducer;
import org.coolsugar.log.service.LogService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

@Service("logService")
public class LogServiceImpl implements LogService {

    @Autowired
    private LogProducer producer;

    public void sendMsg() {
        producer.sendMsg();
    }
}
