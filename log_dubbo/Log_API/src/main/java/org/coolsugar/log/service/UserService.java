package org.coolsugar.log.service;

import org.coolsugar.log.po.User;

public interface UserService {
    User getUserById(int id); //按id查找user

    void saveUser(User user); //将user保存(存在则更新)

    void sayHi();
}
